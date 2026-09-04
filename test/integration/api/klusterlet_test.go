// Copyright Contributors to the Open Cluster Management project
package api

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	operatorv1 "open-cluster-management.io/api/operator/v1"
)

var _ = Describe("Create Klusterlet API", func() {
	var klusterlet *operatorv1.Klusterlet
	BeforeEach(func() {
		suffix := rand.String(5)
		klusterletName := fmt.Sprintf("cm-%s", suffix)
		klusterlet = &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{},
		}
	})

	Context("Create without nothing set", func() {
		It("should create successfully", func() {
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})

	Context("Create with invalid namespace", func() {
		It("should reject the klusterlet creation", func() {
			klusterlet.Spec.Namespace = "invalid-klusterlet-ns"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).NotTo(BeNil())
		})
	})

	Context("Create with aws auth and invalid arn", func() {
		It("should reject the klusterlet creation", func() {
			klusterlet.Spec.RegistrationConfiguration = &operatorv1.RegistrationConfiguration{
				RegistrationDriver: operatorv1.RegistrationDriver{
					AuthType: "awsirsa",
					AwsIrsa: &operatorv1.AwsIrsa{
						ManagedClusterArn: "arn:aws:bks:us-west-2:123456789012:cluster/managed-cluster1",
						HubClusterArn:     "arn:aws:eks:us-west-2:123456789012:cluster/hub-cluster1",
					},
				},
			}
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).NotTo(BeNil())
		})
	})

	Context("Create with aws auth and valid arn", func() {
		It("should create successfully", func() {
			klusterlet.Spec.RegistrationConfiguration = &operatorv1.RegistrationConfiguration{
				RegistrationDriver: operatorv1.RegistrationDriver{
					AuthType: "awsirsa",
					AwsIrsa: &operatorv1.AwsIrsa{
						ManagedClusterArn: "arn:aws:eks:us-west-2:123456789012:cluster/managed-cluster1",
						HubClusterArn:     "arn:aws:eks:us-west-2:123456789012:cluster/hub-cluster1",
					},
				},
			}
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})
})

var _ = Describe("valid HubApiServerHostAlias", func() {
	var klusterlet *operatorv1.Klusterlet

	BeforeEach(func() {
		suffix := rand.String(5)
		klusterletName := fmt.Sprintf("cm-%s", suffix)
		klusterlet = &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{
				HubApiServerHostAlias: &operatorv1.HubApiServerHostAlias{},
			},
		}
	})

	Context("Empty IPV4 address", func() {
		It("should return err", func() {
			klusterlet.Spec.HubApiServerHostAlias.Hostname = "xxx.yyy.zzz"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Empty hostname", func() {
		It("should return err", func() {
			klusterlet.Spec.HubApiServerHostAlias.IP = "1.2.3.4"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Invalid IPV4 address and hostname", func() {
		It("should return err", func() {
			klusterlet.Spec.HubApiServerHostAlias.IP = "1.2.3.257"
			klusterlet.Spec.HubApiServerHostAlias.Hostname = "xxxyyyzzz"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Valid IPV4 address and hostname", func() {
		It("should create successfully", func() {
			klusterlet.Spec.HubApiServerHostAlias.IP = "1.2.3.4"
			klusterlet.Spec.HubApiServerHostAlias.Hostname = "xxx.yyy.zzz"
			_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
	})
})

var _ = Describe("Klusterlet API test with WorkConfiguration", func() {
	var klusterletName string

	BeforeEach(func() {
		suffix := rand.String(5)
		klusterletName = fmt.Sprintf("cm-%s", suffix)
	})

	It("Create a klusterlet with empty spec", func() {
		klusterlet := &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{},
		}
		klusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())

		Expect(klusterlet.Spec.WorkConfiguration).To(BeNil())
	})

	It("Create a klusterlet with empty work feature gate mode", func() {
		klusterlet := &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{
				WorkConfiguration: &operatorv1.WorkAgentConfiguration{
					FeatureGates: []operatorv1.FeatureGate{
						{
							Feature: "Foo",
						},
					},
				},
			},
		}
		klusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).ToNot(HaveOccurred())
		Expect(klusterlet.Spec.WorkConfiguration.FeatureGates[0].Mode).Should(Equal(operatorv1.FeatureGateModeTypeDisable))
	})

	It("Create a klusterlet with wrong work feature gate mode", func() {
		klusterlet := &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{
				WorkConfiguration: &operatorv1.WorkAgentConfiguration{
					FeatureGates: []operatorv1.FeatureGate{
						{
							Feature: "Foo",
							Mode:    "WrongMode",
						},
					},
				},
			},
		}
		_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).To(HaveOccurred())
	})

	It("Create a klusterlet with right work feature gate mode", func() {
		klusterlet := &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{
				Name: klusterletName,
			},
			Spec: operatorv1.KlusterletSpec{
				WorkConfiguration: &operatorv1.WorkAgentConfiguration{
					FeatureGates: []operatorv1.FeatureGate{
						{
							Feature: "Foo",
							Mode:    "Disable",
						},
						{
							Feature: "Bar",
							Mode:    "Enable",
						},
					},
				},
			},
		}
		_, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
		Expect(err).To(BeNil())
		Expect(klusterlet.Spec.WorkConfiguration.FeatureGates[0].Mode).Should(Equal(operatorv1.FeatureGateModeTypeDisable))
		Expect(klusterlet.Spec.WorkConfiguration.FeatureGates[1].Mode).Should(Equal(operatorv1.FeatureGateModeTypeEnable))
	})
})

var _ = Describe("Klusterlet hosting-cluster report validation", func() {
	var createdNames []string

	createKlusterlet := func(spec operatorv1.KlusterletSpec) (*operatorv1.Klusterlet, error) {
		name := fmt.Sprintf("hosting-report-%s", rand.String(5))
		createdNames = append(createdNames, name)
		return operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), &operatorv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Spec:       spec,
		}, metav1.CreateOptions{})
	}

	BeforeEach(func() {
		createdNames = nil
	})

	AfterEach(func() {
		for _, name := range createdNames {
			err := operatorClient.OperatorV1().Klusterlets().Delete(context.TODO(), name, metav1.DeleteOptions{})
			Expect(apierrors.IsNotFound(err) || err == nil).To(BeTrue())
		}
	})

	It("requires clusterName when reporting in a local mode", func() {
		for _, mode := range []operatorv1.InstallMode{"", operatorv1.InstallModeDefault, operatorv1.InstallModeSingleton} {
			_, err := createKlusterlet(operatorv1.KlusterletSpec{
				DeployOption: operatorv1.KlusterletDeployOption{
					Mode:                 mode,
					ReportHostingCluster: operatorv1.ReportHostingClusterModeEnable,
				},
			})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("spec.clusterName is required"))
		}
	})

	It("requires managementClusterName when reporting in a hosted mode", func() {
		for _, mode := range []operatorv1.InstallMode{operatorv1.InstallModeHosted, operatorv1.InstallModeSingletonHosted} {
			_, err := createKlusterlet(operatorv1.KlusterletSpec{
				DeployOption: operatorv1.KlusterletDeployOption{
					Mode:                 mode,
					ReportHostingCluster: operatorv1.ReportHostingClusterModeEnable,
					Hosted:               &operatorv1.KlusterletHostedConfiguration{},
				},
			})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("spec.deployOption.hosted.managementClusterName is required"))
		}
	})

	It("rejects disabling ClusterClaim while reporting", func() {
		for _, featureGateMode := range []operatorv1.FeatureGateModeType{"", operatorv1.FeatureGateModeTypeDisable} {
			_, err := createKlusterlet(operatorv1.KlusterletSpec{
				ClusterName: "spoke-a",
				DeployOption: operatorv1.KlusterletDeployOption{
					Mode:                 operatorv1.InstallModeDefault,
					ReportHostingCluster: operatorv1.ReportHostingClusterModeEnable,
				},
				RegistrationConfiguration: &operatorv1.RegistrationConfiguration{
					FeatureGates: []operatorv1.FeatureGate{{
						Feature: "ClusterClaim",
						Mode:    featureGateMode,
					}},
				},
			})
			Expect(apierrors.IsInvalid(err)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("ClusterClaim cannot be disabled"))
		}
	})

	It("allows reporting with the default or explicitly enabled ClusterClaim setting", func() {
		for _, featureGates := range [][]operatorv1.FeatureGate{
			nil,
			{{Feature: "ClusterClaim", Mode: operatorv1.FeatureGateModeTypeEnable}},
		} {
			_, err := createKlusterlet(operatorv1.KlusterletSpec{
				ClusterName: "spoke-a",
				DeployOption: operatorv1.KlusterletDeployOption{
					Mode:                 operatorv1.InstallModeDefault,
					ReportHostingCluster: operatorv1.ReportHostingClusterModeEnable,
				},
				RegistrationConfiguration: &operatorv1.RegistrationConfiguration{FeatureGates: featureGates},
			})
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("allows ClusterClaim to be disabled when reporting is disabled", func() {
		_, err := createKlusterlet(operatorv1.KlusterletSpec{
			DeployOption: operatorv1.KlusterletDeployOption{
				Mode:                 operatorv1.InstallModeDefault,
				ReportHostingCluster: operatorv1.ReportHostingClusterModeDisable,
			},
			RegistrationConfiguration: &operatorv1.RegistrationConfiguration{
				FeatureGates: []operatorv1.FeatureGate{{
					Feature: "ClusterClaim",
					Mode:    operatorv1.FeatureGateModeTypeDisable,
				}},
			},
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("rejects an update that removes a required hosting-cluster name", func() {
		klusterlet, err := createKlusterlet(operatorv1.KlusterletSpec{
			ClusterName: "spoke-a",
			DeployOption: operatorv1.KlusterletDeployOption{
				Mode:                 operatorv1.InstallModeDefault,
				ReportHostingCluster: operatorv1.ReportHostingClusterModeEnable,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		klusterlet.Spec.ClusterName = ""
		_, err = operatorClient.OperatorV1().Klusterlets().Update(context.TODO(), klusterlet, metav1.UpdateOptions{})
		Expect(apierrors.IsInvalid(err)).To(BeTrue())
		Expect(err.Error()).To(ContainSubstring("spec.clusterName is required"))
	})
})

var _ = Describe("Klusterlet v1 Enhanced API test", func() {
	var klusterletName string

	BeforeEach(func() {
		suffix := rand.String(5)
		klusterletName = fmt.Sprintf("klusterlet-enhanced-%s", suffix)
	})

	AfterEach(func() {
		err := operatorClient.OperatorV1().Klusterlets().Delete(context.TODO(), klusterletName, metav1.DeleteOptions{})
		if !apierrors.IsForbidden(err) {
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Context("Klusterlet comprehensive configuration validation", func() {
		It("should handle complete configuration with all optional fields", func() {
			klusterlet := &operatorv1.Klusterlet{
				ObjectMeta: metav1.ObjectMeta{
					Name: klusterletName,
				},
				Spec: operatorv1.KlusterletSpec{
					RegistrationImagePullSpec: "quay.io/test/registration:latest",
					WorkImagePullSpec:         "quay.io/test/work:latest",
					ClusterName:               "test-cluster",
					Namespace:                 "open-cluster-management-agent",
					ExternalServerURLs: []operatorv1.ServerURL{
						{
							URL: "https://hub.example.com:6443",
						},
					},
					NodePlacement: operatorv1.NodePlacement{
						NodeSelector: map[string]string{
							"node-role.kubernetes.io/worker": "",
						},
						Tolerations: []v1.Toleration{
							{
								Key:      "node-role.kubernetes.io/worker",
								Operator: v1.TolerationOpExists,
								Effect:   v1.TaintEffectNoSchedule,
							},
						},
					},
					DeployOption: operatorv1.KlusterletDeployOption{
						Mode: operatorv1.InstallModeDefault,
					},
					RegistrationConfiguration: &operatorv1.RegistrationConfiguration{
						FeatureGates: []operatorv1.FeatureGate{
							{
								Feature: "AddonManagement",
								Mode:    operatorv1.FeatureGateModeTypeEnable,
							},
						},
					},
					WorkConfiguration: &operatorv1.WorkAgentConfiguration{
						FeatureGates: []operatorv1.FeatureGate{
							{
								Feature: "ManifestWorkReplicaSet",
								Mode:    operatorv1.FeatureGateModeTypeEnable,
							},
						},
					},
				},
			}

			createdKlusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdKlusterlet.Spec.ClusterName).Should(Equal("test-cluster"))
			Expect(createdKlusterlet.Spec.Namespace).Should(Equal("open-cluster-management-agent"))
			Expect(len(createdKlusterlet.Spec.ExternalServerURLs)).Should(Equal(1))
			Expect(createdKlusterlet.Spec.NodePlacement.NodeSelector["node-role.kubernetes.io/worker"]).Should(Equal(""))
			Expect(len(createdKlusterlet.Spec.NodePlacement.Tolerations)).Should(Equal(1))
		})

		It("should validate hosted mode configuration", func() {
			klusterlet := &operatorv1.Klusterlet{
				ObjectMeta: metav1.ObjectMeta{
					Name: klusterletName,
				},
				Spec: operatorv1.KlusterletSpec{
					DeployOption: operatorv1.KlusterletDeployOption{
						Mode: operatorv1.InstallModeHosted,
					},
				},
			}

			createdKlusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdKlusterlet.Spec.DeployOption.Mode).Should(Equal(operatorv1.InstallModeHosted))
		})

		It("should validate priority class configuration", func() {
			klusterlet := &operatorv1.Klusterlet{
				ObjectMeta: metav1.ObjectMeta{
					Name: klusterletName,
				},
				Spec: operatorv1.KlusterletSpec{
					PriorityClassName: "system-cluster-critical",
				},
			}

			createdKlusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdKlusterlet.Spec.PriorityClassName).Should(Equal("system-cluster-critical"))
		})
	})

	Context("Klusterlet resource requirements", func() {
		It("should handle resource requirements configuration", func() {
			klusterlet := &operatorv1.Klusterlet{
				ObjectMeta: metav1.ObjectMeta{
					Name: klusterletName,
				},
				Spec: operatorv1.KlusterletSpec{
					ResourceRequirement: &operatorv1.ResourceRequirement{
						Type: operatorv1.ResourceQosClassResourceRequirement,
					},
				},
			}

			createdKlusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(createdKlusterlet.Spec.ResourceRequirement.Type).Should(Equal(operatorv1.ResourceQosClassResourceRequirement))
		})
	})

	Context("Klusterlet status updates", func() {
		It("should allow status updates", func() {
			klusterlet := &operatorv1.Klusterlet{
				ObjectMeta: metav1.ObjectMeta{
					Name: klusterletName,
				},
				Spec: operatorv1.KlusterletSpec{},
			}

			createdKlusterlet, err := operatorClient.OperatorV1().Klusterlets().Create(context.TODO(), klusterlet, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())

			// Update status
			createdKlusterlet.Status = operatorv1.KlusterletStatus{
				ObservedGeneration: 1,
				Conditions: []metav1.Condition{
					{
						Type:               "Applied",
						Status:             metav1.ConditionTrue,
						Reason:             "KlusterletDeployed",
						LastTransitionTime: metav1.Now(),
					},
				},
				Generations: []operatorv1.GenerationStatus{
					{
						Group:          "apps",
						Version:        "v1",
						Resource:       "deployments",
						Namespace:      "open-cluster-management-agent",
						Name:           "klusterlet-registration-agent",
						LastGeneration: 1,
					},
				},
				RelatedResources: []operatorv1.RelatedResourceMeta{
					{
						Group:     "apps",
						Version:   "v1",
						Resource:  "deployments",
						Namespace: "open-cluster-management-agent",
						Name:      "klusterlet-registration-agent",
					},
				},
			}

			_, err = operatorClient.OperatorV1().Klusterlets().UpdateStatus(context.TODO(), createdKlusterlet, metav1.UpdateOptions{})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
