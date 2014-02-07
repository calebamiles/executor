package runoncehandler_test

import (
	"errors"
	"github.com/cloudfoundry-incubator/executor/actionrunner/fakeactionrunner"
	. "github.com/cloudfoundry-incubator/executor/runoncehandler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/executor/taskregistry/faketaskregistry"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/fakebbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/vito/gordon/fake_gordon"
)

var _ = Describe("RunOnceHandler", func() {
	var (
		handler *RunOnceHandler

		bbs              *fakebbs.FakeExecutorBBS
		runOnce          models.RunOnce
		fakeTaskRegistry *faketaskregistry.FakeTaskRegistry
		gordon           *fake_gordon.FakeGordon
		actionRunner     *fakeactionrunner.FakeActionRunner
	)

	BeforeEach(func() {
		bbs = fakebbs.NewFakeExecutorBBS()
		gordon = fake_gordon.New()
		actionRunner = fakeactionrunner.New()

		fakeTaskRegistry = faketaskregistry.New()

		runOnce = models.RunOnce{
			Guid: "totally-unique",
			Actions: []models.ExecutorAction{
				{models.RunAction{"sudo reboot"}},
			},
		}

		handler = New(bbs, gordon, fakeTaskRegistry, actionRunner)
	})

	Describe("Handling a RunOnce", func() {
		JustBeforeEach(func() {
			handler.RunOnce(runOnce, "executor-id")
		})

		Context("When there are enough resources to claim the RunOnce", func() {
			It("should allocate resources", func() {
				Ω(fakeTaskRegistry.RegisteredRunOnces).Should(ContainElement(runOnce))
			})

			Context("When the RunOnce can be claimed", func() {
				It("claims it", func() {
					Ω(bbs.ClaimedRunOnce.Guid).Should(Equal(runOnce.Guid))
					Ω(bbs.ClaimedRunOnce.ExecutorID).Should(Equal("executor-id"))
				})

				Context("When a container can be made", func() {
					It("should make the container", func() {
						Ω(gordon.CreatedHandles()).Should(HaveLen(1))
					})

					Context("When the RunOnce can be put into the starting state", func() {
						It("should mark the RunOnce as started", func() {
							Ω(bbs.StartedRunOnce.Guid).Should(Equal(runOnce.Guid))
						})

						It("should start running the actions", func() {
							Ω(actionRunner.ContainerHandle).Should(Equal(gordon.CreatedHandles()[0]))
							Ω(actionRunner.Actions).Should(Equal(runOnce.Actions))
						})

						Context("when the RunOnce actions succeed", func() {
							It("should deallocate resources and mark the RunOnce as completed (succesfully)", func() {
								Ω(bbs.CompletedRunOnce.Guid).Should(Equal(runOnce.Guid))
								Ω(bbs.CompletedRunOnce.Failed).Should(BeFalse())

								Ω(fakeTaskRegistry.UnregisteredRunOnces).Should(ContainElement(runOnce))
								Ω(gordon.DestroyedHandles()).Should(HaveLen(1))
							})
						})

						Context("when the RunOnce actions fail", func() {
							BeforeEach(func() {
								actionRunner.RunError = errors.New("Asplosions!")
							})

							It("should deallocate resources and mark the RunOnce as completed (unsuccesfully)", func() {
								Ω(bbs.CompletedRunOnce.Guid).Should(Equal(runOnce.Guid))
								Ω(bbs.CompletedRunOnce.Failed).Should(BeTrue())
								Ω(bbs.CompletedRunOnce.FailureReason).Should(Equal("Asplosions!"))

								Ω(fakeTaskRegistry.UnregisteredRunOnces).Should(ContainElement(runOnce))
								Ω(gordon.DestroyedHandles()).Should(HaveLen(1))
							})
						})
					})

					Context("When the RunOnce fails to be put into the starting state", func() {
						BeforeEach(func() {
							bbs.StartRunOnceErr = errors.New("bam!")
						})

						It("should destroy the container and deallocate resources", func() {
							Ω(gordon.DestroyedHandles()).Should(HaveLen(1))
							Ω(gordon.DestroyedHandles()).Should(Equal(gordon.CreatedHandles()))

							Ω(fakeTaskRegistry.UnregisteredRunOnces).Should(ContainElement(runOnce))
						})
					})
				})

				Context("when a container cannot be made", func() {
					BeforeEach(func() {
						gordon.CreateError = errors.New("No container for you")
					})

					It("does not create a starting RunOnce and it deallocates resources", func() {
						Ω(bbs.StartedRunOnce).Should(BeZero())

						Ω(fakeTaskRegistry.UnregisteredRunOnces).Should(ContainElement(runOnce))
					})
				})
			})

			Context("When the RunOnce cannot be claimed", func() {
				BeforeEach(func() {
					bbs.ClaimRunOnceErr = errors.New("bam!")
				})

				It("does not start the run once, create a container, or set aside resources for it", func() {
					Ω(bbs.StartedRunOnce).Should(BeZero())

					Ω(gordon.CreatedHandles()).Should(BeEmpty())
					Ω(fakeTaskRegistry.UnregisteredRunOnces).Should(ContainElement(runOnce))
				})
			})
		})

		Context("When there are not enough resources to claim the RunOnce", func() {
			BeforeEach(func() {
				fakeTaskRegistry.AddRunOnceErr = errors.New("No room at the inn!")
			})

			It("should not claim the run once or reserve resources", func() {
				Ω(fakeTaskRegistry.RegisteredRunOnces).ShouldNot(ContainElement(runOnce))
				Ω(bbs.ClaimedRunOnce).Should(BeZero())
			})
		})
	})
})
