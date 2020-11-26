package mosaic5g

// TODO make (Requests: corev1.ResourceList) and (Limits: corev1.ResourceList) optional
import (
	"context"
	"fmt"
	"reflect"
	"time"

	Err "errors"

	"mosaic5g/m5g-operator/internal/util"
	mosaic5gv1alpha1 "mosaic5g/m5g-operator/pkg/apis/mosaic5g/v1alpha1"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_mosaic5g")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Mosaic5g Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMosaic5g{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mosaic5g-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Mosaic5g
	err = c.Watch(&source.Kind{Type: &mosaic5gv1alpha1.Mosaic5g{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Mosaic5g
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mosaic5gv1alpha1.Mosaic5g{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMosaic5g{}

// ReconcileMosaic5g reconciles a Mosaic5g object
type ReconcileMosaic5g struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

var databaseName string

var mysqlDatabase *appsv1.Deployment
var mysqlDatabaseDeployment *appsv1.Deployment
var mysqlDatabaseService *v1.Service

var cassandraDatabase *appsv1.Deployment
var cassandraDatabaseDeployment *appsv1.Deployment
var cassandraDatabaseService *v1.Service

// Core network v1
var cnV1 *appsv1.Deployment
var coreNetworkDeploymentV1 *appsv1.Deployment
var coreNetworkServiceV1 *v1.Service

// Core network v2
var cnV2 *appsv1.Deployment
var coreNetworkDeploymentV2 *appsv1.Deployment
var coreNetworkServiceV2 *v1.Service

// oai-hss v1
var hssV1 *appsv1.Deployment
var hssDeploymentV1 *appsv1.Deployment
var hssServiceV1 *v1.Service

// oai-hss v2
var hssV2 *appsv1.Deployment
var hssDeploymentV2 *appsv1.Deployment
var hssServiceV2 *v1.Service

// oai-spgw v1
var spgwV1 *appsv1.Deployment
var spgwDeploymentV1 *appsv1.Deployment
var spgwServiceV1 *v1.Service

// oai-spgwc v2
var spgwcV2 *appsv1.Deployment
var spgwcDeploymentV2 *appsv1.Deployment
var spgwcServiceV2 *v1.Service

// oai-spgwu v2
var spgwuV2 *appsv1.Deployment
var spgwuDeploymentV2 *appsv1.Deployment
var spgwuServiceV2 *v1.Service

// oai-mme v1
var mmeV1 *appsv1.Deployment
var mmeDeploymentV1 *appsv1.Deployment
var mmeServiceV1 *v1.Service

// oai-mme v2
var mmeV2 *appsv1.Deployment
var mmeDeploymentV2 *appsv1.Deployment
var mmeServiceV2 *v1.Service

// oai-ran v1
var ran *appsv1.Deployment
var ranDeployment *appsv1.Deployment
var ranService *v1.Service

// Reconcile reads that state of the cluster for a Mosaic5g object and makes changes based on the state read
// and what is in the Mosaic5g.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
// How to reconcile Mosaic5g:
// 1. Create MySQL, OAI-CN and OAI-RAN in order
// 2. If the configuration changed, restart all OAI PODs
func (r *ReconcileMosaic5g) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// time.Now().Format("2006-01-02 15:04:05"),
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Mosaic5g")

	// Fetch the Mosaic5g instance
	instance := &mosaic5gv1alpha1.Mosaic5g{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue

			//Delete unused ConfigMap
			conf := r.genConfigMap(instance)
			conf.Namespace = "default"
			err = r.client.Delete(context.TODO(), conf)
			if err != nil {
				reqLogger.Error(err, "Failed to delete ")
			}
			reqLogger.Info("Mosaic5g resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Mosaic5g")
		return reconcile.Result{}, err
	}

	new := r.genConfigMap(instance)

	config := &v1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: new.GetName(), Namespace: instance.Namespace}, config)
	if err != nil && errors.IsNotFound(err) {
		// Create a configmap for cn and ran
		reqLogger.Info("Creating a new ConfigMap")
		conf := r.genConfigMap(instance)
		reqLogger.Info("conf", "content", conf)
		err = r.client.Create(context.TODO(), conf)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ConfigMap")
		}
	} else if err != nil {
		reqLogger.Error(err, "Get ConfigMap failed")
		return reconcile.Result{}, err
	}

	/* START the new change from HERE*/
	/*==============================================================================================*/

	/***************************** Create the deployment and service of database *****************************/
	// var databaseName string
	// var databaseDeployment *appsv1.Deployment
	// var databaseService *v1.Service
	//////////
	// var mysqlDatabaseDeployment *appsv1.Deployment
	// var mysqlDatabaseService *v1.Service

	// var cassandraDatabaseDeployment *appsv1.Deployment
	// var cassandraDatabaseService *v1.Service

	//////////

	databaseName = instance.Spec.Database[0].K8sServiceName
	if instance.Spec.Database[0].DatabaseType == "mysql" {
		mysqlDatabase = &appsv1.Deployment{}
		mysqlDatabaseDeployment = r.deploymentForMySQL(instance)
		mysqlDatabaseService = r.genMySQLService(instance)

		// Define a new Database deployment
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: mysqlDatabaseDeployment.GetName(), Namespace: instance.Namespace}, mysqlDatabase)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", mysqlDatabaseDeployment.Namespace, "Deployment.Name", mysqlDatabaseDeployment.Name)
			err = r.client.Create(context.TODO(), mysqlDatabaseDeployment)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", mysqlDatabaseDeployment.Namespace, "Deployment.Name", mysqlDatabaseDeployment.Name)
				return reconcile.Result{}, err
			}

			// Define a new mysql service
			err = r.client.Create(context.TODO(), mysqlDatabaseService)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", mysqlDatabaseService.Namespace, "Service.Name", mysqlDatabaseService.Name)
				return reconcile.Result{}, err
			}

			// Deployment created successfully - return and requeue
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, databaseName, " Failed to get Deployment")
			return reconcile.Result{}, err
		}
	} else {
		cassandraDatabase = &appsv1.Deployment{}
		cassandraDatabaseDeployment = r.deploymentForCassandra(instance)
		cassandraDatabaseService = r.genCassandraService(instance)

		// Define a new Database deployment
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: cassandraDatabaseDeployment.GetName(), Namespace: instance.Namespace}, cassandraDatabase)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", cassandraDatabaseDeployment.Namespace, "Deployment.Name", cassandraDatabaseDeployment.Name)
			err = r.client.Create(context.TODO(), cassandraDatabaseDeployment)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", cassandraDatabaseDeployment.Namespace, "Deployment.Name", cassandraDatabaseDeployment.Name)
				return reconcile.Result{}, err
			}

			// Define a new mysql service
			err = r.client.Create(context.TODO(), cassandraDatabaseService)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", cassandraDatabaseService.Namespace, "Service.Name", cassandraDatabaseService.Name)
				return reconcile.Result{}, err
			}

			// Deployment created successfully - return and requeue
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, databaseName, " Failed to get Deployment")
			return reconcile.Result{}, err
		}
	}
	/***************************** Create the deployment and service of core networks v1 (all-in-one mode) if exist *****************************/

	// cnV1 := &appsv1.Deployment{}
	if len(instance.Spec.OaiCn.V1) >= 1 {
		cnV1 = &appsv1.Deployment{}
		coreNetworkDeploymentV1 = r.deploymentForCnV1(instance)
		coreNetworkServiceV1 = r.genCnV1Service(instance)
		// Check if the oai-cn deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: coreNetworkDeploymentV1.GetName(), Namespace: instance.Namespace}, cnV1)
		if err != nil && errors.IsNotFound(err) {
			if instance.Spec.Database[0].DatabaseType == "mysql" {
				if mysqlDatabase.Status.ReadyReplicas == 0 {
					return reconcile.Result{Requeue: true}, Err.New("No " + databaseName + " POD is ready")
				}
			} else {
				if cassandraDatabase.Status.ReadyReplicas == 0 {
					return reconcile.Result{Requeue: true}, Err.New("No " + databaseName + " POD is ready")
				}
			}
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", coreNetworkDeploymentV1.Namespace, "Deployment.Name", coreNetworkDeploymentV1.Name)
			err = r.client.Create(context.TODO(), coreNetworkDeploymentV1)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", coreNetworkDeploymentV1.Namespace, "Deployment.Name", coreNetworkDeploymentV1.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("30s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "CN Failed to get Deployment")
			return reconcile.Result{}, err
		}

		// Create an oaicn service
		service := &v1.Service{}
		// Check if the oai-cn service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: coreNetworkServiceV1.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), coreNetworkServiceV1)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", coreNetworkServiceV1.Namespace, "Service.Name", coreNetworkServiceV1.Name)
				return reconcile.Result{}, err
			}
		}
	}
	/***************************** Create the deployment and service of core networks v2 (all-in-one mode) if exist *****************************/
	// cnV2 := &appsv1.Deployment{}
	if len(instance.Spec.OaiCn.V2) >= 1 {
		cnV2 = &appsv1.Deployment{}
		coreNetworkDeploymentV2 = r.deploymentForCnV2(instance)
		coreNetworkServiceV2 = r.genCnV2Service(instance)
		// Check if the oai-cn deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: coreNetworkDeploymentV2.GetName(), Namespace: instance.Namespace}, cnV2)
		if err != nil && errors.IsNotFound(err) {
			if instance.Spec.Database[0].DatabaseType == "mysql" {
				if mysqlDatabase.Status.ReadyReplicas == 0 {
					return reconcile.Result{Requeue: true}, Err.New("No " + databaseName + " POD is ready")
				}
			} else {
				if cassandraDatabase.Status.ReadyReplicas == 0 {
					return reconcile.Result{Requeue: true}, Err.New("No " + databaseName + " POD is ready")
				}
			}

			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", coreNetworkDeploymentV2.Namespace, "Deployment.Name", coreNetworkDeploymentV2.Name)
			err = r.client.Create(context.TODO(), coreNetworkDeploymentV2)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", coreNetworkDeploymentV2.Namespace, "Deployment.Name", coreNetworkDeploymentV2.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("30s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "CN Failed to get Deployment")
			return reconcile.Result{}, err
		}

		// Create an oaicn service
		service := &v1.Service{}
		// Check if the oai-cn service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: coreNetworkServiceV2.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), coreNetworkServiceV2)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", coreNetworkServiceV2.Namespace, "Service.Name", coreNetworkServiceV2.Name)
				return reconcile.Result{}, err
			}
		}
	}
	/***************************** Create the deployment and service of oai-hss v1 if exist *****************************/
	// hssV1 := &appsv1.Deployment{}
	if len(instance.Spec.OaiHss.V1) >= 1 {
		hssV1 = &appsv1.Deployment{}
		hssDeploymentV1 = r.deploymentForHssV1(instance)
		hssServiceV1 = r.genHssV1Service(instance)
		// Check if the oai-hss deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: hssDeploymentV1.GetName(), Namespace: instance.Namespace}, hssV1)
		if err != nil && errors.IsNotFound(err) {
			if instance.Spec.Database[0].DatabaseType == "mysql" {
				if mysqlDatabase.Status.ReadyReplicas == 0 {
					return reconcile.Result{Requeue: true}, Err.New("No " + databaseName + " POD is ready")
				}
			} else {
				if cassandraDatabase.Status.ReadyReplicas == 0 {
					return reconcile.Result{Requeue: true}, Err.New("No " + databaseName + " POD is ready")
				}
			}
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", hssDeploymentV1.Namespace, "Deployment.Name", hssDeploymentV1.Name)
			err = r.client.Create(context.TODO(), hssDeploymentV1)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", hssDeploymentV1.Namespace, "Deployment.Name", hssDeploymentV1.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "HSS Failed to get Deployment")
			return reconcile.Result{}, err
		}

		//time.Sleep(15 * time.Second)
		// Create an oaihss service
		service := &v1.Service{}
		// Check if the oai-hss service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: hssServiceV1.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), hssServiceV1)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of HSS", "Service.Namespace", hssServiceV1.Namespace, "Service.Name", hssServiceV1.Name)
				return reconcile.Result{}, err
			}
		}
	}
	/***************************** Create the deployment and service of oai-hss v2 if exist *****************************/
	// hssV2 := &appsv1.Deployment{}
	if len(instance.Spec.OaiHss.V2) >= 1 {
		hssV2 = &appsv1.Deployment{}
		hssDeploymentV2 = r.deploymentForHssV2(instance)
		hssServiceV2 = r.genHssV2Service(instance)
		// Check if the oai-hss deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: hssDeploymentV2.GetName(), Namespace: instance.Namespace}, hssV2)
		if err != nil && errors.IsNotFound(err) {
			if instance.Spec.Database[0].DatabaseType == "mysql" {
				if mysqlDatabase.Status.ReadyReplicas == 0 {
					return reconcile.Result{Requeue: true}, Err.New("No " + databaseName + " POD is ready")
				}
			} else {
				if cassandraDatabase.Status.ReadyReplicas == 0 {
					return reconcile.Result{Requeue: true}, Err.New("No " + databaseName + " POD is ready")
				}
			}
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", hssDeploymentV2.Namespace, "Deployment.Name", hssDeploymentV2.Name)
			err = r.client.Create(context.TODO(), hssDeploymentV2)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", hssDeploymentV2.Namespace, "Deployment.Name", hssDeploymentV2.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "HSS Failed to get Deployment")
			return reconcile.Result{}, err
		}

		//time.Sleep(15 * time.Second)
		// Create an oaihss service
		service := &v1.Service{}
		// Check if the oai-hss service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: hssServiceV2.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), hssServiceV2)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of HSS", "Service.Namespace", hssServiceV2.Namespace, "Service.Name", hssServiceV2.Name)
				return reconcile.Result{}, err
			}
		}
	}

	/***************************** Create the deployment and service of oai-spgw v1 if exist *****************************/
	// spgwV1 := &appsv1.Deployment{}
	if len(instance.Spec.OaiSpgw.V1) >= 1 {
		spgwV1 = &appsv1.Deployment{}
		spgwDeploymentV1 = r.deploymentForSpgwV1(instance)
		spgwServiceV1 = r.genSpgwV1Service(instance)

		// Check if the oai-mme deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: spgwDeploymentV1.GetName(), Namespace: instance.Namespace}, spgwV1)
		if err != nil && errors.IsNotFound(err) {

			if hssV1.Status.ReadyReplicas == 0 {
				return reconcile.Result{Requeue: true}, Err.New("No " + instance.Spec.OaiMme.V1[0].K8sDeploymentName + " POD is ready")
			}
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", spgwDeploymentV1.Namespace, "Deployment.Name", spgwDeploymentV1.Name)
			err = r.client.Create(context.TODO(), spgwDeploymentV1)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", spgwDeploymentV1.Namespace, "Deployment.Name", spgwDeploymentV1.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "MME Failed to get Deployment")
			return reconcile.Result{}, err
		}

		//time.Sleep(5 * time.Second)
		service := &v1.Service{}
		// Check if the oai-spgw service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: spgwServiceV1.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), spgwServiceV1)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of SPGW", "Service.Namespace", spgwServiceV1.Namespace, "Service.Name", spgwServiceV1.Name)
				return reconcile.Result{}, err
			}
		}
		time.Sleep(15 * time.Second)

	}
	/***************************** Create the deployment and service of oai-spgwc v2 if exist *****************************/
	// spgwcV2 := &appsv1.Deployment{}
	if len(instance.Spec.OaiSpgwc.V2) >= 1 {
		spgwcV2 = &appsv1.Deployment{}
		spgwcDeploymentV2 = r.deploymentForSpgwcV2(instance)
		spgwcServiceV2 = r.genSpgwcV2Service(instance)

		// oai-cn v2
		// Check if the oai-mme deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: spgwcDeploymentV2.GetName(), Namespace: instance.Namespace}, spgwcV2)
		if err != nil && errors.IsNotFound(err) {
			if hssV2.Status.ReadyReplicas == 0 {
				return reconcile.Result{Requeue: true}, Err.New("No mme POD is ready")
			}
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", spgwcDeploymentV2.Namespace, "Deployment.Name", spgwcDeploymentV2.Name)
			err = r.client.Create(context.TODO(), spgwcDeploymentV2)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", spgwcDeploymentV2.Namespace, "Deployment.Name", spgwcDeploymentV2.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "MME Failed to get Deployment")
			return reconcile.Result{}, err
		}
		service := &v1.Service{}
		// Check if the oai-spgwc service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: spgwcServiceV2.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), spgwcServiceV2)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of SPGWC", "Service.Namespace", spgwcServiceV2.Namespace, "Service.Name", spgwcServiceV2.Name)
				return reconcile.Result{}, err
			}
		}
		time.Sleep(15 * time.Second)

	}
	/***************************** Create the deployment and service of oai-spgwu v2 if exist *****************************/
	// spgwuV2 := &appsv1.Deployment{}
	if len(instance.Spec.OaiSpgwu.V2) >= 1 {
		spgwuV2 = &appsv1.Deployment{}
		spgwuDeploymentV2 = r.deploymentForSpgwuV2(instance)
		spgwuServiceV2 = r.genSpgwuV2Service(instance)

		// Check if the oai-mme deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: spgwuDeploymentV2.GetName(), Namespace: instance.Namespace}, spgwuV2)
		if err != nil && errors.IsNotFound(err) {
			if hssV2.Status.ReadyReplicas == 0 {
				return reconcile.Result{Requeue: true}, Err.New("No mme POD is ready")
			}
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", spgwuDeploymentV2.Namespace, "Deployment.Name", spgwuDeploymentV2.Name)
			err = r.client.Create(context.TODO(), spgwuDeploymentV2)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", spgwuDeploymentV2.Namespace, "Deployment.Name", spgwuDeploymentV2.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "MME Failed to get Deployment")
			return reconcile.Result{}, err
		}

		service := &v1.Service{}
		// Check if the oai-spgwu service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: spgwuServiceV2.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), spgwuServiceV2)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of SPGWU", "Service.Namespace", spgwuServiceV2.Namespace, "Service.Name", spgwuServiceV2.Name)
				return reconcile.Result{}, err
			}
		}
		time.Sleep(15 * time.Second)

	}
	/***************************** Create the deployment and service of oai-mme v1 if exist *****************************/
	// mmeV1 := &appsv1.Deployment{}
	if len(instance.Spec.OaiMme.V1) >= 1 {
		mmeV1 = &appsv1.Deployment{}
		mmeDeploymentV1 = r.deploymentForMmeV1(instance)
		mmeServiceV1 = r.genMmeV1Service(instance)

		// Creat an oaimme deployment
		// Check if the oai-mme deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: mmeDeploymentV1.GetName(), Namespace: instance.Namespace}, mmeV1)
		if err != nil && errors.IsNotFound(err) {

			if spgwV1.Status.ReadyReplicas == 0 {
				return reconcile.Result{Requeue: true}, Err.New("No " + instance.Spec.OaiSpgw.V1[0].K8sDeploymentName + " POD is ready")
			}
			// switch instance.Spec.Mosaic5gSnapVersion {
			// case "v2":
			// 	if (spgwcV2.Status.ReadyReplicas == 0) && (spgwuV2.Status.ReadyReplicas == 0) {
			// 		return reconcile.Result{Requeue: true}, Err.New("No neither " + instance.Spec.OaiSpgwc.V2[0].K8sDeploymentName + " POD nor " + instance.Spec.OaiSpgwu.V2[0].K8sDeploymentName + " POD are ready")
			// 	}
			// 	if spgwcV2.Status.ReadyReplicas == 0 {
			// 		return reconcile.Result{Requeue: true}, Err.New("No " + instance.Spec.OaiSpgwc.V2[0].K8sDeploymentName + " POD is ready")
			// 	}
			// 	if spgwuV2.Status.ReadyReplicas == 0 {
			// 		return reconcile.Result{Requeue: true}, Err.New("No " + instance.Spec.OaiSpgwu.V2[0].K8sDeploymentName + " POD is ready")
			// 	}

			// default:
			// 	if spgwV1.Status.ReadyReplicas == 0 {
			// 		return reconcile.Result{Requeue: true}, Err.New("No " + instance.Spec.OaiSpgw.V1[0].K8sDeploymentName + " POD is ready")
			// 	}
			// }
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", mmeDeploymentV1.Namespace, "Deployment.Name", mmeDeploymentV1.Name)
			err = r.client.Create(context.TODO(), mmeDeploymentV1)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", mmeDeploymentV1.Namespace, "Deployment.Name", mmeDeploymentV1.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "MME Failed to get Deployment")
			return reconcile.Result{}, err
		}

		service := &v1.Service{}
		// Create an oaimme service
		// Check if the oai-mme service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: mmeServiceV1.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), mmeServiceV1)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of MME", "Service.Namespace", mmeServiceV1.Namespace, "Service.Name", mmeServiceV1.Name)
				return reconcile.Result{}, err
			}
		}
		time.Sleep(15 * time.Second)
	}
	/***************************** Create the deployment and service of oai-mme v2 if exist *****************************/
	// mmeV2 := &appsv1.Deployment{}
	if len(instance.Spec.OaiMme.V2) >= 1 {
		mmeV2 = &appsv1.Deployment{}
		mmeDeploymentV2 = r.deploymentForMmeV2(instance)
		mmeServiceV2 = r.genMmeV2Service(instance)

		// Creat an oaimme deployment
		// Check if the oai-mme deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: mmeDeploymentV2.GetName(), Namespace: instance.Namespace}, mmeV2)
		if err != nil && errors.IsNotFound(err) {

			if (spgwcV2.Status.ReadyReplicas == 0) && (spgwuV2.Status.ReadyReplicas == 0) {
				return reconcile.Result{Requeue: true}, Err.New("No neither " + instance.Spec.OaiSpgwc.V2[0].K8sDeploymentName + " POD nor " + instance.Spec.OaiSpgwu.V2[0].K8sDeploymentName + " POD are ready")
			}
			if spgwcV2.Status.ReadyReplicas == 0 {
				return reconcile.Result{Requeue: true}, Err.New("No " + instance.Spec.OaiSpgwc.V2[0].K8sDeploymentName + " POD is ready")
			}
			if spgwuV2.Status.ReadyReplicas == 0 {
				return reconcile.Result{Requeue: true}, Err.New("No " + instance.Spec.OaiSpgwu.V2[0].K8sDeploymentName + " POD is ready")
			}

			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", mmeDeploymentV2.Namespace, "Deployment.Name", mmeDeploymentV2.Name)
			err = r.client.Create(context.TODO(), mmeDeploymentV2)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", mmeDeploymentV2.Namespace, "Deployment.Name", mmeDeploymentV2.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "MME Failed to get Deployment")
			return reconcile.Result{}, err
		}

		service := &v1.Service{}
		// Create an oaimme service
		// Check if the oai-mme service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: mmeServiceV2.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), mmeServiceV2)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of MME", "Service.Namespace", mmeServiceV2.Namespace, "Service.Name", mmeServiceV2.Name)
				return reconcile.Result{}, err
			}
		}
		time.Sleep(15 * time.Second)
	}
	/***************************** Create the deployment and service of oai-ran if exist *****************************/
	// ran := &appsv1.Deployment{}
	if len(instance.Spec.OaiEnb) >= 1 {
		ran = &appsv1.Deployment{}
		ranDeployment = r.deploymentForRAN(instance)
		ranService = r.genRanService(instance)

		//time.Sleep(15 * time.Second)
		// Create an oairan deployment
		// Check if the oai-ran deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: ranDeployment.GetName(), Namespace: instance.Namespace}, ran)
		if err != nil && errors.IsNotFound(err) {
			// wait for oai-cn V1 if the deployment is all-in-one v1
			if len(instance.Spec.OaiCn.V1) >= 1 {
				if cnV1.Status.ReadyReplicas == 0 {
					d, _ := time.ParseDuration("10s")
					return reconcile.Result{Requeue: true, RequeueAfter: d}, Err.New("No " + instance.Spec.OaiCn.V1[0].K8sDeploymentName + " POD is ready, 10 seconds backoff")
				}
			}
			// wait for oai-cn V2 if the deployment is all-in-one v2
			if len(instance.Spec.OaiCn.V2) >= 1 {
				if cnV2.Status.ReadyReplicas == 0 {
					d, _ := time.ParseDuration("10s")
					return reconcile.Result{Requeue: true, RequeueAfter: d}, Err.New("No " + instance.Spec.OaiCn.V2[0].K8sDeploymentName + " POD is ready, 10 seconds backoff")
				}
			}

			// wait for oai-mme V1 if the deployment is oai-mme v1
			if len(instance.Spec.OaiMme.V1) >= 1 {
				if mmeV1.Status.ReadyReplicas == 0 {
					d, _ := time.ParseDuration("10s")
					return reconcile.Result{Requeue: true, RequeueAfter: d}, Err.New("No " + instance.Spec.OaiMme.V1[0].K8sDeploymentName + " POD is ready, 10 seconds backoff")
				}
			}

			// wait for oai-mme V2 if the deployment is oai-mme v2
			if len(instance.Spec.OaiMme.V2) >= 1 {
				if mmeV2.Status.ReadyReplicas == 0 {
					d, _ := time.ParseDuration("10s")
					return reconcile.Result{Requeue: true, RequeueAfter: d}, Err.New("No " + instance.Spec.OaiMme.V2[0].K8sDeploymentName + " POD is ready, 10 seconds backoff")
				}
			}

			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", ranDeployment.Namespace, "Deployment.Name", ranDeployment.Name)
			err = r.client.Create(context.TODO(), ranDeployment)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", ranDeployment.Namespace, "Deployment.Name", ranDeployment.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "RAN Failed to get Deployment")
			return reconcile.Result{}, err
		}

		// Create an oairan service
		service := &v1.Service{}
		// Check if the oai-cn service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: ranService.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), ranService)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", ranService.Namespace, "Service.Name", ranService.Name)
				return reconcile.Result{}, err
			}
		}
	}

	/* Ensure the deployment size is the same as the spec */
	// Ensure the deployment size of oai-cn v1 (if exist) is the same as the spec
	if len(instance.Spec.OaiCn.V1) >= 1 {
		size := instance.Spec.OaiCn.V1[0].OaiCnSize
		if *cnV1.Spec.Replicas != size {
			cnV1.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), cnV1)

			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", cnV1.Namespace, "Deployment.Name", cnV1.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Ensure the deployment size of oai-cn v2 (if exist) is the same as the spec
	if len(instance.Spec.OaiCn.V2) >= 1 {
		size := instance.Spec.OaiCn.V2[0].OaiCnSize
		if *cnV2.Spec.Replicas != size {
			cnV2.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), cnV2)

			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", cnV2.Namespace, "Deployment.Name", cnV2.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Ensure the deployment size of oai-hss v1 (if exist) is the same as the spec
	if len(instance.Spec.OaiHss.V1) >= 1 {
		size := instance.Spec.OaiHss.V1[0].OaiHssSize
		if *hssV1.Spec.Replicas != size {
			hssV1.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), hssV1)

			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", hssV1.Namespace, "Deployment.Name", hssV1.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Ensure the deployment size of oai-hss v2 (if exist) is the same as the spec
	if len(instance.Spec.OaiHss.V2) >= 1 {
		size := instance.Spec.OaiHss.V2[0].OaiHssSize
		if *hssV2.Spec.Replicas != size {
			hssV2.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), hssV2)

			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", hssV2.Namespace, "Deployment.Name", hssV2.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Ensure the deployment size of oai-mme v1 (if exist) is the same as the spec
	if len(instance.Spec.OaiMme.V1) >= 1 {
		size := instance.Spec.OaiMme.V1[0].OaiMmeSize
		if *mmeV1.Spec.Replicas != size {
			mmeV1.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), mmeV1)

			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", mmeV1.Namespace, "Deployment.Name", mmeV1.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Ensure the deployment size of oai-mme v2 (if exist) is the same as the spec
	if len(instance.Spec.OaiMme.V2) >= 1 {
		size := instance.Spec.OaiMme.V2[0].OaiMmeSize
		if *mmeV2.Spec.Replicas != size {
			mmeV2.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), mmeV2)

			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", mmeV2.Namespace, "Deployment.Name", mmeV2.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Ensure the deployment size of oai-spgw v1 (if exist) is the same as the spec
	if len(instance.Spec.OaiSpgw.V1) >= 1 {
		size := instance.Spec.OaiSpgw.V1[0].OaiSpgwSize
		if *spgwV1.Spec.Replicas != size {
			spgwV1.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), spgwV1)

			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", spgwV1.Namespace, "Deployment.Name", spgwV1.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Ensure the deployment size of oai-spgwc v1 (if exist) is the same as the spec
	if len(instance.Spec.OaiSpgwc.V2) >= 1 {
		size := instance.Spec.OaiSpgwc.V2[0].OaiSpgwcSize
		if *spgwcV2.Spec.Replicas != size {
			spgwcV2.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), spgwcV2)

			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", spgwcV2.Namespace, "Deployment.Name", spgwcV2.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Ensure the deployment size of oai-spgwu v1 (if exist) is the same as the spec
	if len(instance.Spec.OaiSpgwu.V2) >= 1 {
		size := instance.Spec.OaiSpgwu.V2[0].OaiSpgwuSize
		if *spgwuV2.Spec.Replicas != size {
			spgwuV2.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), spgwuV2)

			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", spgwuV2.Namespace, "Deployment.Name", spgwuV2.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// Update the Mosaic5g status with the pod names
	// List the pods for this instance's deployment
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(util.LabelsForMosaic5g(instance.GetName()))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods", "Mosaic5g.Namespace", instance.Namespace, "Mosaic5g.Name", instance.Name)
		return reconcile.Result{}, err
	}
	podNames := util.GetPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Mosaic5g status")
			return reconcile.Result{}, err
		}
	}

	/*==============================================================================================*/
	/* END the new change from HERE*/

	// Check configmap is fine or not. If it's changed, update ConfigMap and restart cn ran
	if err == nil {
		if reflect.DeepEqual(new.Data, config.Data) {
			reqLogger.Info("newConf equals config")
		} else {
			reqLogger.Info("newConf does not equals config")
			reqLogger.Info("Update ConfigMap, deployments and services")
			/*=========================================================*/
			// newconfYaml
			newconfYaml := mosaic5gv1alpha1.Mosaic5gSpec{}
			err := yaml.Unmarshal([]byte(new.Data["conf.yaml"]), &newconfYaml)
			if err != nil {
				reqLogger.Info(err.Error())
			} else {
				fmt.Println(newconfYaml)
			}
			// currentconfYaml
			currentconfYaml := mosaic5gv1alpha1.Mosaic5gSpec{}
			err = yaml.Unmarshal([]byte(config.Data["conf.yaml"]), &currentconfYaml)
			if err != nil {
				reqLogger.Info(err.Error())
			} else {
				fmt.Println(currentconfYaml)
			}

			/*=========================================================*/
			err = r.client.Update(context.TODO(), new)
			// ================================ Database ================================ //
			// Current and new types of database
			currentDatabaseType := currentconfYaml.Database[0].DatabaseType
			newDatabaseType := newconfYaml.Database[0].DatabaseType
			if currentDatabaseType != newDatabaseType {
				if currentDatabaseType == "mysql" {
					size := instance.Spec.Database[0].DatabaseSize
					size = 0
					mysqlDatabase.Spec.Replicas = &size
					err = r.client.Update(context.TODO(), mysqlDatabase)
					if err != nil {
						reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", mysqlDatabase.Namespace, "Deployment.Name", mysqlDatabase.Name)
					}
					fmt.Printf("mysqlDatabaseDeployment.Name=:%v \t mysqlDatabaseService.Name=:%v\n", mysqlDatabaseDeployment.Name, mysqlDatabaseService.Name)
				} else {
					size := instance.Spec.Database[0].DatabaseSize
					size = 0
					cassandraDatabase.Spec.Replicas = &size
					err = r.client.Update(context.TODO(), cassandraDatabase)
					if err != nil {
						reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", cassandraDatabase.Namespace, "Deployment.Name", cassandraDatabase.Name)
					}
					fmt.Printf("cassandraDatabaseDeployment.Name=:%v \t cassandraDatabaseService.Name=:%v\n", cassandraDatabaseDeployment.Name, cassandraDatabaseService.Name)
				}

			}

			// // update database if exist
			// if len(instance.Spec.Database) >= 1 {
			// 	if instance.Spec.Database[0].DatabaseType == "mysql" {
			// 		err = r.client.Delete(context.TODO(), mysqlDatabaseDeployment)
			// 		err = r.client.Delete(context.TODO(), mysqlDatabaseService)

			// 	} else {
			// 		err = r.client.Delete(context.TODO(), cassandraDatabaseDeployment)
			// 		err = r.client.Delete(context.TODO(), cassandraDatabaseService)
			// 	}
			// }

			// ================================ oai-cn v1 ================================ //
			if (len(currentconfYaml.OaiCn.V1) >= 1) && (len(newconfYaml.OaiCn.V1) == 0) {
				// Delet the core network v1 from the previous deployment
				size := currentconfYaml.OaiCn.V1[0].OaiCnSize
				size = 0
				cnV1.Spec.Replicas = &size
				err = r.client.Update(context.TODO(), cnV1)
				if err != nil {
					reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", cnV1.Namespace, "Deployment.Name", cnV1.Name)
				}

			}
			// ================================ oai-cn v2 ================================ //
			if (len(currentconfYaml.OaiCn.V2) >= 1) && (len(newconfYaml.OaiCn.V2) == 0) {
				// Delet the core network v2 from the previous deployment
				size := currentconfYaml.OaiCn.V2[0].OaiCnSize
				size = 0
				cnV2.Spec.Replicas = &size
				err = r.client.Update(context.TODO(), cnV2)
				if err != nil {
					reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", cnV2.Namespace, "Deployment.Name", cnV2.Name)
				}

			}

			// ================================ oai-hss v1 ================================ //
			if (len(currentconfYaml.OaiHss.V1) >= 1) && (len(newconfYaml.OaiHss.V1) == 0) {
				// Delet the core network v1 from the previous deployment
				size := currentconfYaml.OaiHss.V1[0].OaiHssSize
				size = 0
				hssV1.Spec.Replicas = &size
				err = r.client.Update(context.TODO(), hssV1)
				if err != nil {
					reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", hssV1.Namespace, "Deployment.Name", hssV1.Name)
				}
			}
			// ================================ oai-hss v2 ================================ //
			if (len(currentconfYaml.OaiHss.V2) >= 1) && (len(newconfYaml.OaiHss.V2) == 0) {
				// Delet the core network v2 from the previous deployment
				size := currentconfYaml.OaiHss.V2[0].OaiHssSize
				size = 0
				hssV2.Spec.Replicas = &size
				err = r.client.Update(context.TODO(), hssV2)
				if err != nil {
					reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", hssV2.Namespace, "Deployment.Name", hssV2.Name)
				}
			}
			// ================================ oai-mme v1 ================================ //
			if (len(currentconfYaml.OaiMme.V1) >= 1) && (len(newconfYaml.OaiMme.V1) == 0) {
				// Delet the core network v1 from the previous deployment
				size := currentconfYaml.OaiMme.V1[0].OaiMmeSize
				size = 0
				mmeV1.Spec.Replicas = &size
				err = r.client.Update(context.TODO(), mmeV1)
				if err != nil {
					reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", mmeV1.Namespace, "Deployment.Name", mmeV1.Name)
				}
			}
			// ================================ oai-mme v2 ================================ //
			if (len(currentconfYaml.OaiMme.V2) >= 1) && (len(newconfYaml.OaiMme.V2) == 0) {
				// Delet the core network v2 from the previous deployment
				size := currentconfYaml.OaiMme.V2[0].OaiMmeSize
				size = 0
				mmeV2.Spec.Replicas = &size
				err = r.client.Update(context.TODO(), mmeV2)
				if err != nil {
					reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", mmeV2.Namespace, "Deployment.Name", mmeV2.Name)
				}
			}

			// ================================ oai-spgw v1 ================================ //
			if (len(currentconfYaml.OaiSpgw.V1) >= 1) && (len(newconfYaml.OaiSpgw.V1) == 0) {
				// Delet the core network v1 from the previous deployment
				size := currentconfYaml.OaiSpgw.V1[0].OaiSpgwSize
				size = 0
				spgwV1.Spec.Replicas = &size
				err = r.client.Update(context.TODO(), spgwV1)
				if err != nil {
					reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", spgwV1.Namespace, "Deployment.Name", spgwV1.Name)
				}
			}

			// ================================ oai-spgwc v2 ================================ //
			if (len(currentconfYaml.OaiSpgwc.V2) >= 1) && (len(newconfYaml.OaiSpgwc.V2) == 0) {
				// Delet the core network v2 from the previous deployment
				size := currentconfYaml.OaiSpgwc.V2[0].OaiSpgwcSize
				size = 0
				spgwcV2.Spec.Replicas = &size
				err = r.client.Update(context.TODO(), spgwcV2)
				if err != nil {
					reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", spgwcV2.Namespace, "Deployment.Name", spgwcV2.Name)
				}
			}

			// ================================ oai-spgwu v2 ================================ //
			if (len(currentconfYaml.OaiSpgwu.V2) >= 1) && (len(newconfYaml.OaiSpgwu.V2) == 0) {
				// Delet the core network v2 from the previous deployment
				size := currentconfYaml.OaiSpgwu.V2[0].OaiSpgwuSize
				size = 0
				spgwuV2.Spec.Replicas = &size
				err = r.client.Update(context.TODO(), spgwuV2)
				if err != nil {
					reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", spgwuV2.Namespace, "Deployment.Name", spgwuV2.Name)
				}
			}

			// update oai-cn v1 if exist
			if len(instance.Spec.OaiCn.V1) >= 1 {
				err = r.client.Delete(context.TODO(), coreNetworkDeploymentV1)
				err = r.client.Delete(context.TODO(), coreNetworkServiceV1)
			}

			// update oai-cn v1 if exist
			if len(instance.Spec.OaiCn.V2) >= 1 {
				err = r.client.Delete(context.TODO(), coreNetworkDeploymentV2)
				err = r.client.Delete(context.TODO(), coreNetworkServiceV2)
			}

			// update oai-hss v1 if exist
			if len(instance.Spec.OaiHss.V1) >= 1 {
				err = r.client.Delete(context.TODO(), hssDeploymentV1)
				err = r.client.Delete(context.TODO(), hssServiceV1)
			}

			// update oai-hss v2 if exist
			if len(instance.Spec.OaiHss.V2) >= 1 {
				err = r.client.Delete(context.TODO(), hssDeploymentV2)
				err = r.client.Delete(context.TODO(), hssServiceV2)
			}

			// update oai-mme v1 if exist
			if len(instance.Spec.OaiMme.V1) >= 1 {
				err = r.client.Delete(context.TODO(), mmeDeploymentV1)
				err = r.client.Delete(context.TODO(), mmeServiceV1)
			}

			// update oai-mme v2 if exist
			if len(instance.Spec.OaiMme.V2) >= 1 {
				err = r.client.Delete(context.TODO(), mmeDeploymentV2)
				err = r.client.Delete(context.TODO(), mmeServiceV2)
			}

			// update oai-spgw v1 if exist
			if len(instance.Spec.OaiSpgw.V1) >= 1 {
				err = r.client.Delete(context.TODO(), spgwDeploymentV1)
				err = r.client.Delete(context.TODO(), spgwServiceV1)
			}

			// update oai-spgwc v2 if exist
			if len(instance.Spec.OaiSpgwc.V2) >= 1 {
				err = r.client.Delete(context.TODO(), spgwcDeploymentV2)
				err = r.client.Delete(context.TODO(), spgwcServiceV2)
			}

			// update oai-spgwu v2 if exist
			if len(instance.Spec.OaiSpgwu.V2) >= 1 {
				err = r.client.Delete(context.TODO(), spgwuDeploymentV2)
				err = r.client.Delete(context.TODO(), spgwuServiceV2)
			}

			// update oai-spgwu v2 if exist
			if len(instance.Spec.OaiEnb) >= 1 {
				err = r.client.Delete(context.TODO(), ranDeployment)
				err = r.client.Delete(context.TODO(), ranService)
			}
			/////////////////////////////////////////////////////////////////////////////////////
			// /////////////////////////////////////////////////////////////////////////////////////
			// // update oai-cn v1 if exist
			// if len(instance.Spec.OaiCn.V1) >= 1 {
			// 	err = r.client.Update(context.TODO(), coreNetworkDeploymentV1)
			// 	err = r.client.Update(context.TODO(), coreNetworkServiceV1)
			// }

			// // update oai-cn v1 if exist
			// if len(instance.Spec.OaiCn.V2) >= 1 {
			// 	err = r.client.Update(context.TODO(), coreNetworkDeploymentV2)
			// 	err = r.client.Update(context.TODO(), coreNetworkServiceV2)
			// }

			// // update oai-hss v1 if exist
			// if len(instance.Spec.OaiHss.V1) >= 1 {
			// 	err = r.client.Update(context.TODO(), hssDeploymentV1)
			// 	err = r.client.Update(context.TODO(), hssServiceV1)
			// }

			// // update oai-hss v2 if exist
			// if len(instance.Spec.OaiHss.V2) >= 1 {
			// 	err = r.client.Update(context.TODO(), hssDeploymentV2)
			// 	err = r.client.Update(context.TODO(), hssServiceV2)
			// }

			// // update oai-mme v1 if exist
			// if len(instance.Spec.OaiMme.V1) >= 1 {
			// 	err = r.client.Update(context.TODO(), mmeDeploymentV1)
			// 	err = r.client.Update(context.TODO(), mmeServiceV1)
			// }

			// // update oai-mme v2 if exist
			// if len(instance.Spec.OaiMme.V2) >= 1 {
			// 	err = r.client.Update(context.TODO(), mmeDeploymentV2)
			// 	err = r.client.Update(context.TODO(), mmeServiceV2)
			// }

			// // update oai-spgw v1 if exist
			// if len(instance.Spec.OaiSpgw.V1) >= 1 {
			// 	err = r.client.Update(context.TODO(), spgwDeploymentV1)
			// 	err = r.client.Update(context.TODO(), spgwServiceV1)
			// }

			// // update oai-spgwc v2 if exist
			// if len(instance.Spec.OaiSpgwc.V2) >= 1 {
			// 	err = r.client.Update(context.TODO(), spgwcDeploymentV2)
			// 	err = r.client.Update(context.TODO(), spgwcServiceV2)
			// }

			// // update oai-spgwu v2 if exist
			// if len(instance.Spec.OaiSpgwu.V2) >= 1 {
			// 	err = r.client.Update(context.TODO(), spgwuDeploymentV2)
			// 	err = r.client.Update(context.TODO(), spgwuServiceV2)
			// }

			// // update oai-spgwu v2 if exist
			// if len(instance.Spec.OaiEnb) >= 1 {
			// 	err = r.client.Update(context.TODO(), ranDeployment)
			// 	err = r.client.Update(context.TODO(), ranService)
			// }
			// /////////////////////////////////////////////////////////////////////////////////////
			// Spec updated - return and requeue
			d, _ := time.ParseDuration("10s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		}

	}
	// Everything is fine, Reconcile ends
	return reconcile.Result{}, nil
}

// deploymentForHssV1 returns a HSS Network Deployment object
func (r *ReconcileMosaic5g) deploymentForHssV1(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.OaiHss.V1[0].K8sDeploymentName

	fmt.Println("m.Spec.OaiHss.V1[0]=", m.Spec.OaiHss.V1[0])
	//ls := util.LabelsForMosaic5g(m.Name + hssName)
	replicas := m.Spec.OaiHss.V1[0].OaiHssSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiHss.V1[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiHss.V1[0].K8sLabelSelector[i].Key, m.Spec.OaiHss.V1[0].K8sLabelSelector[i].Value)
		labels[m.Spec.OaiHss.V1[0].K8sLabelSelector[i].Key] = m.Spec.OaiHss.V1[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiHss.V1[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiHss.V1[0].K8sNodeSelector[i].Key, m.Spec.OaiHss.V1[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.OaiHss.V1[0].K8sNodeSelector[i].Key] = m.Spec.OaiHss.V1[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiHss.V1[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.GetName() + "-" + deploymentName,
			// Namespace:   m.Namespace,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					// NodeSelector: map[string]string{
					// 	"usrp": "true"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiHss.V1[0].OaiHssImage,
						Name:            m.Spec.OaiHss.V1[0].K8sDeploymentName,
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiHss.V1[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiHss.V1[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiHss.V1[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiHss.V1[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "cgroup",
								ReadOnly:  true,
								MountPath: "/sys/fs/cgroup/",
							}, {
								Name:      "module",
								ReadOnly:  true,
								MountPath: "/lib/modules/",
							}, {
								Name:      "mosaic5g-config",
								MountPath: "/root/config",
							}},
						/* TODO add the configuration of the ports to the configuration of the deployed in the yaml file to be dynamic.
						Hints: use the folloiwng method to add the ports
						dep.Spec.Template.Spec.Containers[0].Ports[0].Name = "mosaic5g-cn"
						dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 80
						...
						*/
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Name:          "mosaic5g-cn",
							}, {
								ContainerPort: 2152,
								Name:          "hss-1",
							}, {
								ContainerPort: 3868,
								Name:          "hss-2",
							}, {
								ContainerPort: 5868,
								Name:          "hss-3",
							}, {
								ContainerPort: 2123,
								Name:          "hss-4",
							}, {
								ContainerPort: 3870,
								Name:          "hss-5",
							}, {
								ContainerPort: 5870,
								Name:          "hss-6",
							}},
					}},
					Affinity: util.GenAffinity("cn"),
					Volumes: []corev1.Volume{
						{
							Name: "cgroup",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/sys/fs/cgroup/",
									Type: util.NewHostPathType("Directory"),
								},
							}},
						{
							Name: "module",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/lib/modules/",
									Type: util.NewHostPathType("Directory"),
								},
							}}, {
							Name: "mosaic5g-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
								},
							}},
					},
				},
			},
		},
	}
	////////////////////
	if (m.Spec.OaiHss.V1[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.OaiHss.V1[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiHss.V1[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiHss.V1[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.OaiHss.V1[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.OaiHss.V1[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiHss.V1[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiHss.V1[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForHssV2 returns a HSS Network Deployment object
func (r *ReconcileMosaic5g) deploymentForHssV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.OaiHss.V2[0].K8sDeploymentName

	//ls := util.LabelsForMosaic5g(m.Name + hssName)
	replicas := m.Spec.OaiHss.V2[0].OaiHssSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiHss.V2[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiHss.V2[0].K8sLabelSelector[i].Key, m.Spec.OaiHss.V2[0].K8sLabelSelector[i].Value)
		labels[m.Spec.OaiHss.V2[0].K8sLabelSelector[i].Key] = m.Spec.OaiHss.V2[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiHss.V2[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiHss.V2[0].K8sNodeSelector[i].Key, m.Spec.OaiHss.V2[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.OaiHss.V2[0].K8sNodeSelector[i].Key] = m.Spec.OaiHss.V2[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiHss.V2[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.GetName() + "-" + deploymentName,
			// Namespace:   m.Namespace,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					// NodeSelector: map[string]string{
					// 	"usrp": "true"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiHss.V2[0].OaiHssImage,
						Name:            m.Spec.OaiHss.V2[0].K8sDeploymentName,
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiHss.V2[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiHss.V2[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiHss.V2[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiHss.V2[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "cgroup",
								ReadOnly:  true,
								MountPath: "/sys/fs/cgroup/",
							}, {
								Name:      "module",
								ReadOnly:  true,
								MountPath: "/lib/modules/",
							}, {
								Name:      "mosaic5g-config",
								MountPath: "/root/config",
							}},
						/* TODO add the configuration of the ports to the configuration of the deployed in the yaml file to be dynamic.
						Hints: use the folloiwng method to add the ports
						dep.Spec.Template.Spec.Containers[0].Ports[0].Name = "mosaic5g-cn"
						dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 80
						...
						*/
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Name:          "mosaic5g-cn",
							}, {
								ContainerPort: 2152,
								Name:          "hss-1",
							}, {
								ContainerPort: 3868,
								Name:          "hss-2",
							}, {
								ContainerPort: 5868,
								Name:          "hss-3",
							}, {
								ContainerPort: 2123,
								Name:          "hss-4",
							}, {
								ContainerPort: 3870,
								Name:          "hss-5",
							}, {
								ContainerPort: 5870,
								Name:          "hss-6",
							}},
					}},
					Affinity: util.GenAffinity("cn"),
					Volumes: []corev1.Volume{
						{
							Name: "cgroup",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/sys/fs/cgroup/",
									Type: util.NewHostPathType("Directory"),
								},
							}},
						{
							Name: "module",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/lib/modules/",
									Type: util.NewHostPathType("Directory"),
								},
							}}, {
							Name: "mosaic5g-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
								},
							}},
					},
				},
			},
		},
	}
	////////////////////
	if (m.Spec.OaiHss.V2[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.OaiHss.V2[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiHss.V2[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiHss.V2[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.OaiHss.V2[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.OaiHss.V2[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiHss.V2[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiHss.V2[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////

	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForMmeV1 returns a MME Network Deployment object
func (r *ReconcileMosaic5g) deploymentForMmeV1(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.OaiMme.V1[0].K8sDeploymentName

	//ls := util.LabelsForMosaic5g(m.Name + deploymentName)
	replicas := m.Spec.OaiMme.V1[0].OaiMmeSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiMme.V1[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiMme.V1[0].K8sLabelSelector[i].Key, m.Spec.OaiMme.V1[0].K8sLabelSelector[i].Value)
		labels[m.Spec.OaiMme.V1[0].K8sLabelSelector[i].Key] = m.Spec.OaiMme.V1[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiMme.V1[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiMme.V1[0].K8sNodeSelector[i].Key, m.Spec.OaiMme.V1[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.OaiMme.V1[0].K8sNodeSelector[i].Key] = m.Spec.OaiMme.V1[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiMme.V1[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.GetName() + "-" + deploymentName,
			// Namespace:   m.Namespace,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					Hostname:     "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiMme.V1[0].OaiMmeImage,
						Name:            m.Spec.OaiMme.V1[0].K8sDeploymentName,
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiMme.V1[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiMme.V1[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiMme.V1[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiMme.V1[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "cgroup",
							ReadOnly:  true,
							MountPath: "/sys/fs/cgroup/",
						}, {
							Name:      "module",
							ReadOnly:  true,
							MountPath: "/lib/modules/",
						}, {
							Name:      "mosaic5g-config",
							MountPath: "/root/config",
						}},
						/* TODO add the configuration of the ports to the configuration of the deployed in the yaml file to be dynamic.
						Hints: use the folloiwng method to add the ports
						dep.Spec.Template.Spec.Containers[0].Ports[0].Name = "mosaic5g-cn"
						dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 80
						...
						*/
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Name:          "mosaic5g-cn",
							}, {
								ContainerPort: 2152,
								Name:          "mme-1",
							}, {
								ContainerPort: 3868,
								Name:          "mme-2",
							}, {
								ContainerPort: 5868,
								Name:          "mme-3",
							}, {
								ContainerPort: 2123,
								Name:          "mme-4",
							}, {
								ContainerPort: 3870,
								Name:          "mme-5",
							}, {
								ContainerPort: 5870,
								Name:          "mme-6",
							}},
					}},
					Affinity: util.GenAffinity("cn"),
					Volumes: []corev1.Volume{{
						Name: "cgroup",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/sys/fs/cgroup/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "module",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/lib/modules/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "mosaic5g-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
							},
						}},
					},
				},
			},
		},
	}
	////////////////////
	if (m.Spec.OaiMme.V1[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.OaiMme.V1[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiMme.V1[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiMme.V1[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.OaiMme.V1[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.OaiMme.V1[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiMme.V1[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiMme.V1[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForMmeV2 returns a MME Network Deployment object
func (r *ReconcileMosaic5g) deploymentForMmeV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.OaiMme.V2[0].K8sDeploymentName

	//ls := util.LabelsForMosaic5g(m.Name + deploymentName)
	replicas := m.Spec.OaiMme.V2[0].OaiMmeSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiMme.V2[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiMme.V2[0].K8sLabelSelector[i].Key, m.Spec.OaiMme.V2[0].K8sLabelSelector[i].Value)
		labels[m.Spec.OaiMme.V2[0].K8sLabelSelector[i].Key] = m.Spec.OaiMme.V2[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiMme.V2[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiMme.V2[0].K8sNodeSelector[i].Key, m.Spec.OaiMme.V2[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.OaiMme.V2[0].K8sNodeSelector[i].Key] = m.Spec.OaiMme.V2[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiMme.V2[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.GetName() + "-" + deploymentName,
			// Namespace:   m.Namespace,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					Hostname:     "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiMme.V2[0].OaiMmeImage,
						Name:            m.Spec.OaiMme.V2[0].K8sDeploymentName,
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiMme.V2[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiMme.V2[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiMme.V2[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiMme.V2[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "cgroup",
							ReadOnly:  true,
							MountPath: "/sys/fs/cgroup/",
						}, {
							Name:      "module",
							ReadOnly:  true,
							MountPath: "/lib/modules/",
						}, {
							Name:      "mosaic5g-config",
							MountPath: "/root/config",
						}},
						/* TODO add the configuration of the ports to the configuration of the deployed in the yaml file to be dynamic.
						Hints: use the folloiwng method to add the ports
						dep.Spec.Template.Spec.Containers[0].Ports[0].Name = "mosaic5g-cn"
						dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 80
						...
						*/
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Name:          "mosaic5g-cn",
							}, {
								ContainerPort: 2152,
								Name:          "mme-1",
							}, {
								ContainerPort: 3868,
								Name:          "mme-2",
							}, {
								ContainerPort: 5868,
								Name:          "mme-3",
							}, {
								ContainerPort: 2123,
								Name:          "mme-4",
							}, {
								ContainerPort: 3870,
								Name:          "mme-5",
							}, {
								ContainerPort: 5870,
								Name:          "mme-6",
							}},
					}},
					Affinity: util.GenAffinity("cn"),
					Volumes: []corev1.Volume{{
						Name: "cgroup",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/sys/fs/cgroup/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "module",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/lib/modules/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "mosaic5g-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
							},
						}},
					},
				},
			},
		},
	}
	////////////////////
	if (m.Spec.OaiMme.V2[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.OaiMme.V2[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiMme.V2[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiMme.V2[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.OaiMme.V2[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.OaiMme.V2[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiMme.V2[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiMme.V2[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForSpgwV1 returns a SPGW Network Deployment object
func (r *ReconcileMosaic5g) deploymentForSpgwV1(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.OaiSpgw.V1[0].K8sDeploymentName

	//ls := util.LabelsForMosaic5g(m.Name + spgwName)
	replicas := m.Spec.OaiSpgw.V1[0].OaiSpgwSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiSpgw.V1[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiSpgw.V1[0].K8sLabelSelector[i].Key, m.Spec.OaiSpgw.V1[0].K8sLabelSelector[i].Value)
		labels[m.Spec.OaiSpgw.V1[0].K8sLabelSelector[i].Key] = m.Spec.OaiSpgw.V1[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiSpgw.V1[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiSpgw.V1[0].K8sNodeSelector[i].Key, m.Spec.OaiSpgw.V1[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.OaiSpgw.V1[0].K8sNodeSelector[i].Key] = m.Spec.OaiSpgw.V1[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiSpgw.V1[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.GetName() + "-" + deploymentName,
			// Namespace:   m.Namespace,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiSpgw.V1[0].OaiSpgwImage,
						Name:            m.Spec.OaiSpgw.V1[0].K8sDeploymentName,
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgw.V1[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgw.V1[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgw.V1[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgw.V1[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "cgroup",
							ReadOnly:  true,
							MountPath: "/sys/fs/cgroup/",
						}, {
							Name:      "module",
							ReadOnly:  true,
							MountPath: "/lib/modules/",
						}, {
							Name:      "mosaic5g-config",
							MountPath: "/root/config",
						}},
						/* TODO add the configuration of the ports to the configuration of the deployed in the yaml file to be dynamic.
						Hints: use the folloiwng method to add the ports
						dep.Spec.Template.Spec.Containers[0].Ports[0].Name = "mosaic5g-cn"
						dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 80
						...
						*/
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Name:          "mosaic5g-cn",
							}, {
								ContainerPort: 2152,
								Name:          "spgw-1",
							}, {
								ContainerPort: 3868,
								Name:          "spgw-2",
							}, {
								ContainerPort: 5868,
								Name:          "spgw-3",
							}, {
								ContainerPort: 2123,
								Name:          "spgw-4",
							}, {
								ContainerPort: 3870,
								Name:          "spgw-5",
							}, {
								ContainerPort: 5870,
								Name:          "spgw-6",
							}},
					}},
					Affinity: util.GenAffinity("cn"),
					Volumes: []corev1.Volume{{
						Name: "cgroup",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/sys/fs/cgroup/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "module",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/lib/modules/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "mosaic5g-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
							},
						}},
					},
				},
			},
		},
	}
	////////////////////
	if (m.Spec.OaiSpgw.V1[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.OaiSpgw.V1[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgw.V1[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgw.V1[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.OaiSpgw.V1[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.OaiSpgw.V1[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgw.V1[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgw.V1[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForSpgwcV2 returns a SPGW Network Deployment object
func (r *ReconcileMosaic5g) deploymentForSpgwcV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.OaiSpgwc.V2[0].K8sDeploymentName

	//ls := util.LabelsForMosaic5g(m.Name + spgwName)
	replicas := m.Spec.OaiSpgwc.V2[0].OaiSpgwcSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiSpgwc.V2[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiSpgwc.V2[0].K8sLabelSelector[i].Key, m.Spec.OaiSpgwc.V2[0].K8sLabelSelector[i].Value)
		labels[m.Spec.OaiSpgwc.V2[0].K8sLabelSelector[i].Key] = m.Spec.OaiSpgwc.V2[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiSpgwc.V2[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiSpgwc.V2[0].K8sNodeSelector[i].Key, m.Spec.OaiSpgwc.V2[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.OaiSpgwc.V2[0].K8sNodeSelector[i].Key] = m.Spec.OaiSpgwc.V2[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiSpgwc.V2[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.GetName() + "-" + deploymentName,
			// Namespace:   m.Namespace,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiSpgwc.V2[0].OaiSpgwcImage,
						Name:            m.Spec.OaiSpgwc.V2[0].K8sDeploymentName,
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgwc.V2[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgwc.V2[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgwc.V2[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgwc.V2[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "cgroup",
							ReadOnly:  true,
							MountPath: "/sys/fs/cgroup/",
						}, {
							Name:      "module",
							ReadOnly:  true,
							MountPath: "/lib/modules/",
						}, {
							Name:      "mosaic5g-config",
							MountPath: "/root/config",
						}},
						/* TODO add the configuration of the ports to the configuration of the deployed in the yaml file to be dynamic.
						Hints: use the folloiwng method to add the ports
						dep.Spec.Template.Spec.Containers[0].Ports[0].Name = "mosaic5g-cn"
						dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 80
						...
						*/
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Name:          "mosaic5g-cn",
							}, {
								ContainerPort: 2152,
								Name:          "spgw-1",
							}, {
								ContainerPort: 3868,
								Name:          "spgw-2",
							}, {
								ContainerPort: 5868,
								Name:          "spgw-3",
							}, {
								ContainerPort: 2123,
								Name:          "spgw-4",
							}, {
								ContainerPort: 3870,
								Name:          "spgw-5",
							}, {
								ContainerPort: 5870,
								Name:          "spgw-6",
							}},
					}},
					Affinity: util.GenAffinity("cn"),
					Volumes: []corev1.Volume{{
						Name: "cgroup",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/sys/fs/cgroup/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "module",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/lib/modules/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "mosaic5g-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
							},
						}},
					},
				},
			},
		},
	}
	////////////////////
	if (m.Spec.OaiSpgwc.V2[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.OaiSpgwc.V2[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgwc.V2[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgwc.V2[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.OaiSpgwc.V2[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.OaiSpgwc.V2[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgwc.V2[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgwc.V2[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForSpgwuV2 returns a SPGW Network Deployment object
func (r *ReconcileMosaic5g) deploymentForSpgwuV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.OaiSpgwu.V2[0].K8sDeploymentName

	//ls := util.LabelsForMosaic5g(m.Name + spgwName)
	replicas := m.Spec.OaiSpgwu.V2[0].OaiSpgwuSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiSpgwu.V2[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiSpgwu.V2[0].K8sLabelSelector[i].Key, m.Spec.OaiSpgwu.V2[0].K8sLabelSelector[i].Value)
		labels[m.Spec.OaiSpgwu.V2[0].K8sLabelSelector[i].Key] = m.Spec.OaiSpgwu.V2[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiSpgwu.V2[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiSpgwu.V2[0].K8sNodeSelector[i].Key, m.Spec.OaiSpgwu.V2[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.OaiSpgwu.V2[0].K8sNodeSelector[i].Key] = m.Spec.OaiSpgwu.V2[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiSpgwu.V2[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.GetName() + "-" + deploymentName,
			// Namespace:   m.Namespace,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiSpgwu.V2[0].OaiSpgwuImage,
						Name:            m.Spec.OaiSpgwu.V2[0].K8sDeploymentName,
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgwu.V2[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgwu.V2[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgwu.V2[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgwu.V2[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "cgroup",
							ReadOnly:  true,
							MountPath: "/sys/fs/cgroup/",
						}, {
							Name:      "module",
							ReadOnly:  true,
							MountPath: "/lib/modules/",
						}, {
							Name:      "mosaic5g-config",
							MountPath: "/root/config",
						}},
						/* TODO add the configuration of the ports to the configuration of the deployed in the yaml file to be dynamic.
						Hints: use the folloiwng method to add the ports
						dep.Spec.Template.Spec.Containers[0].Ports[0].Name = "mosaic5g-cn"
						dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 80
						...
						*/
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Name:          "mosaic5g-cn",
							}, {
								ContainerPort: 2152,
								Name:          "spgw-1",
							}, {
								ContainerPort: 3868,
								Name:          "spgw-2",
							}, {
								ContainerPort: 5868,
								Name:          "spgw-3",
							}, {
								ContainerPort: 2123,
								Name:          "spgw-4",
							}, {
								ContainerPort: 3870,
								Name:          "spgw-5",
							}, {
								ContainerPort: 5870,
								Name:          "spgw-6",
							}},
					}},
					Affinity: util.GenAffinity("cn"),
					Volumes: []corev1.Volume{{
						Name: "cgroup",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/sys/fs/cgroup/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "module",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/lib/modules/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "mosaic5g-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
							},
						}},
					},
				},
			},
		},
	}
	////////////////////
	if (m.Spec.OaiSpgwu.V2[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.OaiSpgwu.V2[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgwu.V2[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgwu.V2[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.OaiSpgwu.V2[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.OaiSpgwu.V2[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiSpgwu.V2[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiSpgwu.V2[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForCnV1 returns a Core Network Deployment object
func (r *ReconcileMosaic5g) deploymentForCnV1(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.OaiCn.V1[0].K8sDeploymentName
	// ls := util.LabelsForMosaic5g(m.Name + deploymentName)

	replicas := m.Spec.OaiCn.V1[0].OaiCnSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiCn.V1[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiCn.V1[0].K8sLabelSelector[i].Key, m.Spec.OaiCn.V1[0].K8sLabelSelector[i].Value)
		labels[m.Spec.OaiCn.V1[0].K8sLabelSelector[i].Key] = m.Spec.OaiCn.V1[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiCn.V1[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiCn.V1[0].K8sNodeSelector[i].Key, m.Spec.OaiCn.V1[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.OaiCn.V1[0].K8sNodeSelector[i].Key] = m.Spec.OaiCn.V1[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiCn.V1[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.GetName() + "-" + deploymentName,
			Namespace: namespace,
			//Namespace: m.Namespace,
			Labels:      labels,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiCn.V1[0].OaiCnImage,
						Name:            m.Spec.OaiCn.V1[0].K8sDeploymentName,
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiCn.V1[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiCn.V1[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiCn.V1[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiCn.V1[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "cgroup",
							ReadOnly:  true,
							MountPath: "/sys/fs/cgroup/",
						}, {
							Name:      "module",
							ReadOnly:  true,
							MountPath: "/lib/modules/",
						}, {
							Name:      "mosaic5g-config",
							MountPath: "/root/config",
						}},
						/* TODO add the configuration of the ports to the configuration of the deployed in the yaml file to be dynamic.
						Hints: use the folloiwng method to add the ports
						dep.Spec.Template.Spec.Containers[0].Ports[0].Name = "mosaic5g-cn"
						dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 80
						...
						*/
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "mosaic5g-cn",
						}},
					}},
					Affinity: util.GenAffinity("cn"),
					Volumes: []corev1.Volume{{
						Name: "cgroup",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/sys/fs/cgroup/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "module",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/lib/modules/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "mosaic5g-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
							},
						}},
					},
				},
			},
		},
	}
	////////////////////
	if (m.Spec.OaiCn.V1[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.OaiCn.V1[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiCn.V1[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiCn.V1[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.OaiCn.V1[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.OaiCn.V1[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiCn.V1[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiCn.V1[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForCnV2 returns a Core Network Deployment object
func (r *ReconcileMosaic5g) deploymentForCnV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.OaiCn.V2[0].K8sDeploymentName
	// ls := util.LabelsForMosaic5g(m.Name + deploymentName)

	replicas := m.Spec.OaiCn.V2[0].OaiCnSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiCn.V2[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiCn.V2[0].K8sLabelSelector[i].Key, m.Spec.OaiCn.V2[0].K8sLabelSelector[i].Value)
		labels[m.Spec.OaiCn.V2[0].K8sLabelSelector[i].Key] = m.Spec.OaiCn.V2[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiCn.V2[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiCn.V2[0].K8sNodeSelector[i].Key, m.Spec.OaiCn.V2[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.OaiCn.V2[0].K8sNodeSelector[i].Key] = m.Spec.OaiCn.V2[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiCn.V2[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.GetName() + "-" + deploymentName,
			Namespace: namespace,
			//Namespace: m.Namespace,
			Labels:      labels,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiCn.V2[0].OaiCnImage,
						Name:            m.Spec.OaiCn.V2[0].K8sDeploymentName,
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiCn.V2[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiCn.V2[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiCn.V2[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiCn.V2[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "cgroup",
							ReadOnly:  true,
							MountPath: "/sys/fs/cgroup/",
						}, {
							Name:      "module",
							ReadOnly:  true,
							MountPath: "/lib/modules/",
						}, {
							Name:      "mosaic5g-config",
							MountPath: "/root/config",
						}},
						/* TODO add the configuration of the ports to the configuration of the deployed in the yaml file to be dynamic.
						Hints: use the folloiwng method to add the ports
						dep.Spec.Template.Spec.Containers[0].Ports[0].Name = "mosaic5g-cn"
						dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 80
						...
						*/
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "mosaic5g-cn",
						}},
					}},
					Affinity: util.GenAffinity("cn"),
					Volumes: []corev1.Volume{{
						Name: "cgroup",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/sys/fs/cgroup/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "module",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/lib/modules/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "mosaic5g-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
							},
						}},
					},
				},
			},
		},
	}
	////////////////////
	if (m.Spec.OaiCn.V2[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.OaiCn.V2[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiCn.V2[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiCn.V2[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.OaiCn.V2[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.OaiCn.V2[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiCn.V2[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiCn.V2[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForRAN returns a Core Network Deployment object
func (r *ReconcileMosaic5g) deploymentForRAN(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.OaiEnb[0].K8sDeploymentName

	// ls := util.LabelsForMosaic5g(m.Name)
	replicas := m.Spec.OaiEnb[0].OaiEnbSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiEnb[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiEnb[0].K8sLabelSelector[i].Key, m.Spec.OaiEnb[0].K8sLabelSelector[i].Value)
		labels[m.Spec.OaiEnb[0].K8sLabelSelector[i].Key] = m.Spec.OaiEnb[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiEnb[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiEnb[0].K8sNodeSelector[i].Key, m.Spec.OaiEnb[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.OaiEnb[0].K8sNodeSelector[i].Key] = m.Spec.OaiEnb[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiEnb[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.GetName() + "-" + deploymentName,
			// Namespace:   m.Namespace,
			Namespace:   namespace,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					// NodeSelector: map[string]string{
					// 	"usrp": "true"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiEnb[0].OaiEnbImage,
						Name:            m.Spec.OaiEnb[0].K8sDeploymentName,
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiEnb[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiEnb[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiEnb[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.OaiEnb[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "cgroup",
							ReadOnly:  true,
							MountPath: "/sys/fs/cgroup/",
						}, {
							Name:      "module",
							ReadOnly:  true,
							MountPath: "/lib/modules/",
						}, {
							Name:      "usrp",
							ReadOnly:  true,
							MountPath: "/dev/bus/usb/",
						}, {
							Name:      "mosaic5g-config",
							MountPath: "/root/config",
						}},
						/* TODO add the configuration of the ports to the configuration of the deployed in the yaml file to be dynamic.
						Hints: use the folloiwng method to add the ports
						dep.Spec.Template.Spec.Containers[0].Ports[0].Name = "mosaic5g-cn"
						dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 80
						...
						*/
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 80,
								Name:          "mosaic5g-ran",
							}, {
								ContainerPort: 2210,
								Name:          "ran-1",
							}, {
								ContainerPort: 22100,
								Name:          "ran-2",
							}, {
								ContainerPort: 2152,
								Name:          "ran-3",
							}, {
								ContainerPort: 50000,
								Name:          "ran-4",
							}, {
								ContainerPort: 50001,
								Name:          "ran-5",
							}, {
								ContainerPort: 36412,
								Name:          "ran-6",
							}},
					}},
					Affinity: util.GenAffinity("ran"),
					Volumes: []corev1.Volume{{
						Name: "cgroup",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/sys/fs/cgroup/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "module",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/lib/modules/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "mosaic5g-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
							},
						}}, {
						Name: "usrp",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/dev/bus/usb/",
								Type: util.NewHostPathType("Directory"),
							},
						}},
					},
				},
			},
		},
	}
	////////////////////
	if (m.Spec.OaiEnb[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.OaiEnb[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiEnb[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiEnb[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.OaiEnb[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.OaiEnb[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.OaiEnb[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.OaiEnb[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////

	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForMySQL returns a Core Network Deployment object
func (r *ReconcileMosaic5g) deploymentForMySQL(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.Database[0].K8sDeploymentName

	//ls := util.LabelsForMosaic5g(m.Name + cnName)
	replicas := m.Spec.Database[0].DatabaseSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.Database[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.Database[0].K8sLabelSelector[i].Key, m.Spec.Database[0].K8sLabelSelector[i].Value)
		labels[m.Spec.Database[0].K8sLabelSelector[i].Key] = m.Spec.Database[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.Database[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.Database[0].K8sNodeSelector[i].Key, m.Spec.Database[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.Database[0].K8sNodeSelector[i].Key] = m.Spec.Database[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiEnb[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.GetName() + "-" + deploymentName,
			// Namespace: m.Namespace,
			Namespace:   namespace,
			Annotations: Annotations,
			Labels:      labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Containers: []corev1.Container{{
						Image: m.Spec.Database[0].DatabaseImage, //m.Spec.MysqlImage,
						Name:  m.Spec.Database[0].K8sDeploymentName,
						Env: []corev1.EnvVar{
							{Name: "MYSQL_ROOT_PASSWORD", Value: "linux"},
						},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.Database[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.Database[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.Database[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.Database[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
					}},
					Affinity: util.GenAffinity("database"),
				},
			},
		},
	}
	////////////////////
	if (m.Spec.Database[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.Database[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.Database[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.Database[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.Database[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.Database[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.Database[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.Database[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForCassandra returns a Core Network Deployment object
func (r *ReconcileMosaic5g) deploymentForCassandra(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	deploymentName := m.Spec.Database[0].K8sDeploymentName

	//ls := util.LabelsForMosaic5g(m.Name + cnName)
	replicas := m.Spec.Database[0].DatabaseSize
	labels := make(map[string]string)
	for i := 0; i < len(m.Spec.Database[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.Database[0].K8sLabelSelector[i].Key, m.Spec.Database[0].K8sLabelSelector[i].Value)
		labels[m.Spec.Database[0].K8sLabelSelector[i].Key] = m.Spec.Database[0].K8sLabelSelector[i].Value
	}
	nodeSelctors := make(map[string]string)
	for i := 0; i < len(m.Spec.Database[0].K8sNodeSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.Database[0].K8sNodeSelector[i].Key, m.Spec.Database[0].K8sNodeSelector[i].Value)
		nodeSelctors[m.Spec.Database[0].K8sNodeSelector[i].Key] = m.Spec.Database[0].K8sNodeSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiEnb[0].K8sEntityNamespace
	}
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+deploymentName] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.GetName() + "-" + deploymentName,
			// Namespace: m.Namespace,
			Namespace:   namespace,
			Annotations: Annotations,
			Labels:      labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelctors,
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Containers: []corev1.Container{{
						Image: m.Spec.Database[0].DatabaseImage, //m.Spec.MysqlImage,
						Name:  m.Spec.Database[0].K8sDeploymentName,
						Env: []corev1.EnvVar{
							{Name: "CASSANDRA_CLUSTER_NAME", Value: "OAI HSS Cluster"},
							{Name: "CASSANDRA_ENDPOINT_SNITCH", Value: "GossipingPropertyFileSnitch"},
						},
						// Resources: corev1.ResourceRequirements{
						// 	Limits: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.Database[0].K8sPodResources.Limits.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.Database[0].K8sPodResources.Limits.ResourceMemory),
						// 	},
						// 	Requests: corev1.ResourceList{
						// 		corev1.ResourceCPU:    resource.MustParse(m.Spec.Database[0].K8sPodResources.Requests.ResourceCPU),
						// 		corev1.ResourceMemory: resource.MustParse(m.Spec.Database[0].K8sPodResources.Requests.ResourceMemory),
						// 	},
						// },
						// Ports: []corev1.ContainerPort{{
						// 	ContainerPort: 3306,
						// 	Name:          "mysql",
						// }},
					}},
					Affinity: util.GenAffinity("database"),
				},
			},
		},
	}
	////////////////////
	if (m.Spec.Database[0].K8sPodResources.Limits.ResourceCPU != "") && (m.Spec.Database[0].K8sPodResources.Limits.ResourceMemory != "") {
		limits := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.Database[0].K8sPodResources.Limits.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.Database[0].K8sPodResources.Limits.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Limits = limits
	}
	if (m.Spec.Database[0].K8sPodResources.Requests.ResourceCPU != "") && (m.Spec.Database[0].K8sPodResources.Requests.ResourceMemory != "") {
		requests := corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(m.Spec.Database[0].K8sPodResources.Requests.ResourceCPU),
			corev1.ResourceMemory: resource.MustParse(m.Spec.Database[0].K8sPodResources.Requests.ResourceMemory),
		}
		dep.Spec.Template.Spec.Containers[0].Resources.Requests = requests
	}
	////////////////////
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// genConfigMap will generate a configmap from ReconcileMosaic5g's spec
func (r *ReconcileMosaic5g) genConfigMap(m *mosaic5gv1alpha1.Mosaic5g) *v1.ConfigMap {
	genLogger := log.WithValues("Mosaic5g", "genConfigMap")
	// Make specs into map[name][value]
	datas := make(map[string]string)
	d, err := yaml.Marshal(&m.Spec)
	if err != nil {
		log.Error(err, "Marshal fail")
	}
	datas["conf.yaml"] = string(d)
	cm := v1.ConfigMap{
		Data: datas,
	}
	cm.Name = "mosaic5g-config"
	cm.Namespace = m.Namespace
	genLogger.Info("Done")
	return &cm
}

// genConfigMapForOaiEnb will generate a configmap from ReconcileMosaic5g's spec for oaiEnb
func (r *ReconcileMosaic5g) genConfigMapForOaiEnb(m *mosaic5gv1alpha1.Mosaic5g) *v1.ConfigMap {
	genLogger := log.WithValues("Mosaic5g", "genConfigMap")
	// Make specs into map[name][value]
	datas := make(map[string]string)
	d, err := yaml.Marshal(&m.Spec.OaiEnb)
	if err != nil {
		log.Error(err, "Marshal fail")
	}
	datas["conf.yaml"] = string(d)
	cm := v1.ConfigMap{
		Data: datas,
	}
	cm.Name = "mosaic5g-config-oaienb"
	cm.Namespace = m.Namespace
	genLogger.Info("Done")
	return &cm
}

// genConfigMapForFlexran will generate a configmap from ReconcileMosaic5g's spec for flexran
func (r *ReconcileMosaic5g) genConfigMapForFlexran(m *mosaic5gv1alpha1.Mosaic5g) *v1.ConfigMap {
	genLogger := log.WithValues("Mosaic5g", "genConfigMap")
	// Make specs into map[name][value]
	datas := make(map[string]string)
	d, err := yaml.Marshal(&m.Spec.Flexran)
	if err != nil {
		log.Error(err, "Marshal fail")
	}
	datas["conf.yaml"] = string(d)
	cm := v1.ConfigMap{
		Data: datas,
	}
	cm.Name = "mosaic5g-config-flexran"
	cm.Namespace = m.Namespace
	genLogger.Info("Done")
	return &cm
}

// genConfigMapForLlmec will generate a configmap from ReconcileMosaic5g's spec for llmec
func (r *ReconcileMosaic5g) genConfigMapForLlmec(m *mosaic5gv1alpha1.Mosaic5g) *v1.ConfigMap {
	genLogger := log.WithValues("Mosaic5g", "genConfigMap")
	// Make specs into map[name][value]
	datas := make(map[string]string)
	d, err := yaml.Marshal(&m.Spec.LlMec)
	if err != nil {
		log.Error(err, "Marshal fail")
	}
	datas["conf.yaml"] = string(d)
	cm := v1.ConfigMap{
		Data: datas,
	}
	cm.Name = "mosaic5g-config-llmec"
	cm.Namespace = m.Namespace
	genLogger.Info("Done")
	return &cm
}

// genConfigMapForoaiCnV1 will generate a configmap from ReconcileMosaic5g's spec for oaiCn V1
func (r *ReconcileMosaic5g) genConfigMapForoaiCnV1(m *mosaic5gv1alpha1.Mosaic5g) *v1.ConfigMap {
	genLogger := log.WithValues("Mosaic5g", "genConfigMap")
	// Make specs into map[name][value]
	datas := make(map[string]string)
	d, err := yaml.Marshal(&m.Spec.OaiCn)
	if err != nil {
		log.Error(err, "Marshal fail")
	}
	datas["conf.yaml"] = string(d)
	cm := v1.ConfigMap{
		Data: datas,
	}
	cm.Name = "mosaic5g-config-oaicnv1"
	cm.Namespace = m.Namespace
	genLogger.Info("Done")
	return &cm
}

// genCnV1Service will generate a service for oaicn
func (r *ReconcileMosaic5g) genCnV1Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.OaiCn.V1[0].K8sServiceName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiCn.V1[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiCn.V1[0].K8sLabelSelector[i].Key, m.Spec.OaiCn.V1[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.OaiCn.V1[0].K8sLabelSelector[i].Key] = m.Spec.OaiCn.V1[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiCn.V1[0].K8sEntityNamespace
	}
	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{Name: "enb", Port: 2152},
			{Name: "hss-1", Port: 3868},
			{Name: "hss-2", Port: 5868},
			{Name: "mme", Port: 2123},
			{Name: "spgw-1", Port: 3870},
			{Name: "spgw-2", Port: 5870},
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genCnV2Service will generate a service for oaicn
func (r *ReconcileMosaic5g) genCnV2Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.OaiCn.V2[0].K8sServiceName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiCn.V2[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiCn.V2[0].K8sLabelSelector[i].Key, m.Spec.OaiCn.V2[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.OaiCn.V2[0].K8sLabelSelector[i].Key] = m.Spec.OaiCn.V2[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiCn.V2[0].K8sEntityNamespace
	}
	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{Name: "enb", Port: 2152},
			{Name: "hss-1", Port: 3868},
			{Name: "hss-2", Port: 5868},
			{Name: "mme", Port: 2123},
			{Name: "spgw-1", Port: 3870},
			{Name: "spgw-2", Port: 5870},
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genHssV1Service will generate a service for oai-hss v1
func (r *ReconcileMosaic5g) genHssV1Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.OaiHss.V1[0].K8sDeploymentName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiHss.V1[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiHss.V1[0].K8sLabelSelector[i].Key, m.Spec.OaiHss.V1[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.OaiHss.V1[0].K8sLabelSelector[i].Key] = m.Spec.OaiHss.V1[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiHss.V1[0].K8sEntityNamespace
	}

	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		// Ports: []v1.ServicePort{
		// 	{Name: "hss-enb", Port: 2152},
		// 	{Name: "hss-hss-1", Port: 3868},
		// 	{Name: "hss-hss-2", Port: 5868},
		// 	{Name: "hss-mme", Port: 2123},
		// 	{Name: "hss-spgw-1", Port: 3870},
		// 	{Name: "hss-spgw-2", Port: 5870},
		// },
		Ports: []v1.ServicePort{
			{Name: "hss-enb-5", Port: 80, Protocol: v1.ProtocolTCP}, //tcp
			{Name: "hss-enb", Port: 2152},
			{Name: "hss-hss-1", Port: 3868},
			{Name: "hss-hss-2", Port: 5868},
			{Name: "hss-mme", Port: 2123},
			{Name: "hss-spgw-1", Port: 3870},
			{Name: "hss-spgw-2", Port: 5870},
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genHssV2Service will generate a service for oai-hss v1
func (r *ReconcileMosaic5g) genHssV2Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.OaiHss.V2[0].K8sDeploymentName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiHss.V2[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiHss.V2[0].K8sLabelSelector[i].Key, m.Spec.OaiHss.V2[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.OaiHss.V2[0].K8sLabelSelector[i].Key] = m.Spec.OaiHss.V2[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiHss.V2[0].K8sEntityNamespace
	}

	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		// Ports: []v1.ServicePort{
		// 	{Name: "hss-enb", Port: 2152},
		// 	{Name: "hss-hss-1", Port: 3868},
		// 	{Name: "hss-hss-2", Port: 5868},
		// 	{Name: "hss-mme", Port: 2123},
		// 	{Name: "hss-spgw-1", Port: 3870},
		// 	{Name: "hss-spgw-2", Port: 5870},
		// },
		Ports: []v1.ServicePort{
			{Name: "hss-enb-5", Port: 80, Protocol: v1.ProtocolTCP}, //tcp
			{Name: "hss-enb", Port: 2152},
			{Name: "hss-hss-1", Port: 3868},
			{Name: "hss-hss-2", Port: 5868},
			{Name: "hss-mme", Port: 2123},
			{Name: "hss-spgw-1", Port: 3870},
			{Name: "hss-spgw-2", Port: 5870},
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genMmeV1Service will generate a service for oai-mme v1
func (r *ReconcileMosaic5g) genMmeV1Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.OaiMme.V1[0].K8sDeploymentName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiMme.V1[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiMme.V1[0].K8sLabelSelector[i].Key, m.Spec.OaiMme.V1[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.OaiMme.V1[0].K8sLabelSelector[i].Key] = m.Spec.OaiMme.V1[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiMme.V1[0].K8sEntityNamespace
	}

	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		// Ports: []v1.ServicePort{
		// 	{Name: "mme-enb", Port: 2152},
		// 	{Name: "mme-hss-1", Port: 3868},
		// 	{Name: "mme-hss-2", Port: 5868},
		// 	{Name: "mme-mme", Port: 2123},
		// 	{Name: "mme-spgw-1", Port: 3870},
		// 	{Name: "mme-spgw-2", Port: 5870},
		// },
		Ports: []v1.ServicePort{
			{Name: "mme-enb-5", Port: 80, Protocol: v1.ProtocolTCP}, //tcp
			{Name: "mme-enb", Port: 2152},
			{Name: "mme-hss-1", Port: 3868},
			{Name: "mme-hss-2", Port: 5868},
			{Name: "mme-mme", Port: 2123},
			{Name: "mme-spgw-1", Port: 3870},
			{Name: "mme-spgw-2", Port: 5870},
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genMmeV2Service will generate a service for oai-mme v1
func (r *ReconcileMosaic5g) genMmeV2Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.OaiMme.V2[0].K8sDeploymentName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiMme.V2[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiMme.V2[0].K8sLabelSelector[i].Key, m.Spec.OaiMme.V2[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.OaiMme.V2[0].K8sLabelSelector[i].Key] = m.Spec.OaiMme.V2[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiMme.V2[0].K8sEntityNamespace
	}

	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		// Ports: []v1.ServicePort{
		// 	{Name: "mme-enb", Port: 2152},
		// 	{Name: "mme-hss-1", Port: 3868},
		// 	{Name: "mme-hss-2", Port: 5868},
		// 	{Name: "mme-mme", Port: 2123},
		// 	{Name: "mme-spgw-1", Port: 3870},
		// 	{Name: "mme-spgw-2", Port: 5870},
		// },
		Ports: []v1.ServicePort{
			{Name: "mme-enb-5", Port: 80, Protocol: v1.ProtocolTCP}, //tcp
			{Name: "mme-enb", Port: 2152},
			{Name: "mme-hss-1", Port: 3868},
			{Name: "mme-hss-2", Port: 5868},
			{Name: "mme-mme", Port: 2123},
			{Name: "mme-spgw-1", Port: 3870},
			{Name: "mme-spgw-2", Port: 5870},
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genSpgwV1Service will generate a service for oai-spgw v1
func (r *ReconcileMosaic5g) genSpgwV1Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.OaiSpgw.V1[0].K8sDeploymentName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiSpgw.V1[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiSpgw.V1[0].K8sLabelSelector[i].Key, m.Spec.OaiSpgw.V1[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.OaiSpgw.V1[0].K8sLabelSelector[i].Key] = m.Spec.OaiSpgw.V1[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiSpgw.V1[0].K8sEntityNamespace
	}

	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		// Ports: []v1.ServicePort{
		// 	{Name: "spgw-enb", Port: 2152},
		// 	{Name: "spgw-hss-1", Port: 3868},
		// 	{Name: "spgw-hss-2", Port: 5868},
		// 	{Name: "spgw-mme", Port: 2123},
		// 	{Name: "spgw-spgw-1", Port: 3870},
		// 	{Name: "spgw-spgw-2", Port: 5870},
		// },
		Ports: []v1.ServicePort{
			{Name: "spgw-enb-5", Port: 80, Protocol: v1.ProtocolTCP}, //tcp
			{Name: "spgw-enb", Port: 2152},
			{Name: "spgw-hss-1", Port: 3868},
			{Name: "spgw-hss-2", Port: 5868},
			{Name: "spgw-mme", Port: 2123},
			{Name: "spgw-spgw-1", Port: 3870},
			{Name: "spgw-spgw-2", Port: 5870},
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genSpgwcV2Service will generate a service for oai-spgwc v2
func (r *ReconcileMosaic5g) genSpgwcV2Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.OaiSpgwc.V2[0].K8sDeploymentName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiSpgwc.V2[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiSpgwc.V2[0].K8sLabelSelector[i].Key, m.Spec.OaiSpgwc.V2[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.OaiSpgwc.V2[0].K8sLabelSelector[i].Key] = m.Spec.OaiSpgwc.V2[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiSpgwc.V2[0].K8sEntityNamespace
	}

	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		// Ports: []v1.ServicePort{
		// 	{Name: "spgw-enb", Port: 2152},
		// 	{Name: "spgw-hss-1", Port: 3868},
		// 	{Name: "spgw-hss-2", Port: 5868},
		// 	{Name: "spgw-mme", Port: 2123},
		// 	{Name: "spgw-spgw-1", Port: 3870},
		// 	{Name: "spgw-spgw-2", Port: 5870},
		// },
		Ports: []v1.ServicePort{
			{Name: "spgw-enb-5", Port: 80, Protocol: v1.ProtocolTCP}, //tcp
			{Name: "spgw-enb", Port: 2152},
			{Name: "spgw-hss-1", Port: 3868},
			{Name: "spgw-hss-2", Port: 5868},
			{Name: "spgw-mme", Port: 2123},
			{Name: "spgw-spgw-1", Port: 3870},
			{Name: "spgw-spgw-2", Port: 5870},
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genSpgwuV2Service will generate a service for oai-spgwu v2
func (r *ReconcileMosaic5g) genSpgwuV2Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.OaiSpgwu.V2[0].K8sDeploymentName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiSpgwu.V2[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiSpgwu.V2[0].K8sLabelSelector[i].Key, m.Spec.OaiSpgwu.V2[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.OaiSpgwu.V2[0].K8sLabelSelector[i].Key] = m.Spec.OaiSpgwu.V2[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiSpgwu.V2[0].K8sEntityNamespace
	}

	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		// Ports: []v1.ServicePort{
		// 	{Name: "spgw-enb", Port: 2152},
		// 	{Name: "spgw-hss-1", Port: 3868},
		// 	{Name: "spgw-hss-2", Port: 5868},
		// 	{Name: "spgw-mme", Port: 2123},
		// 	{Name: "spgw-spgw-1", Port: 3870},
		// 	{Name: "spgw-spgw-2", Port: 5870},
		// },
		Ports: []v1.ServicePort{
			{Name: "spgw-enb-5", Port: 80, Protocol: v1.ProtocolTCP}, //tcp
			{Name: "spgw-enb", Port: 2152},
			{Name: "spgw-hss-1", Port: 3868},
			{Name: "spgw-hss-2", Port: 5868},
			{Name: "spgw-mme", Port: 2123},
			{Name: "spgw-spgw-1", Port: 3870},
			{Name: "spgw-spgw-2", Port: 5870},
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genRanService will generate a service for oaicn
func (r *ReconcileMosaic5g) genRanService(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.OaiEnb[0].K8sServiceName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.OaiEnb[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.OaiEnb[0].K8sLabelSelector[i].Key, m.Spec.OaiEnb[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.OaiEnb[0].K8sLabelSelector[i].Key] = m.Spec.OaiEnb[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiEnb[0].K8sEntityNamespace
	}

	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{Name: "enb-enb-5", Port: 80, Protocol: v1.ProtocolTCP},    //tcp
			{Name: "enb-enb", Port: 2210, Protocol: v1.ProtocolTCP},    //tcp
			{Name: "enb-enb-1", Port: 22100, Protocol: v1.ProtocolTCP}, //tcp
			{Name: "enb-s1-u", Port: 2152, Protocol: v1.ProtocolUDP},   //udp
			{Name: "enb-enb-3", Port: 50000, Protocol: v1.ProtocolUDP}, //udp
			{Name: "enb-enb-4", Port: 50001, Protocol: v1.ProtocolUDP}, //udp
			{Name: "enb-s1-c", Port: 36412, Protocol: v1.ProtocolTCP},  //tcp
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genMySQLService will generate a service for oaicn
func (r *ReconcileMosaic5g) genMySQLService(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.Database[0].K8sServiceName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.Database[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.Database[0].K8sLabelSelector[i].Key, m.Spec.Database[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.Database[0].K8sLabelSelector[i].Key] = m.Spec.Database[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiEnb[0].K8sEntityNamespace
	}
	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{Name: "mysql-port", Port: 3306},
		},
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genCassandraService will generate a service for mysql
func (r *ReconcileMosaic5g) genCassandraService(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	serviceName := m.Spec.Database[0].K8sServiceName

	var service *v1.Service

	selectMap := make(map[string]string)
	for i := 0; i < len(m.Spec.Database[0].K8sLabelSelector); i++ {
		fmt.Printf("key=:%v \t key=:%v\n", m.Spec.Database[0].K8sLabelSelector[i].Key, m.Spec.Database[0].K8sLabelSelector[i].Value)
		selectMap[m.Spec.Database[0].K8sLabelSelector[i].Key] = m.Spec.Database[0].K8sLabelSelector[i].Value
	}

	namespace := m.Spec.K8sGlobalNamespace
	if namespace == "" {
		namespace = m.Spec.OaiEnb[0].K8sEntityNamespace
	}
	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		// Ports: []v1.ServicePort{
		// 	{Name: "mysql-port", Port: 3306},
		// },
		Selector: selectMap,
		// Type:     "NodePort",
		ClusterIP: "None",
	}
	service.Name = serviceName
	service.Namespace = namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}
