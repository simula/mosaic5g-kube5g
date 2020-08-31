package mosaic5g

import (
	"context"
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

// Reconcile reads that state of the cluster for a Mosaic5g object and makes changes based on the state read
// and what is in the Mosaic5g.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
// How to reconcile Mosaic5g:
// 1. Create CASSANDRA, OAI-CN and OAI-RAN in order
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

	// Define a new CASSANDRA deployment
	cassandradb := &appsv1.Deployment{}
	cassandraDeployment := r.deploymentForCassandra(instance)
	// Check if CASSANDRA deployment already exists, if not create a new one
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: cassandraDeployment.GetName(), Namespace: instance.Namespace}, cassandradb)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", cassandraDeployment.Namespace, "Deployment.Name", cassandraDeployment.Name)
		err = r.client.Create(context.TODO(), cassandraDeployment)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", cassandraDeployment.Namespace, "Deployment.Name", cassandraDeployment.Name)
			return reconcile.Result{}, err
		}
		// Define a new cassandradb service
		cassandraService := r.genCassandraService(instance)
		err = r.client.Create(context.TODO(), cassandraService)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", cassandraService.Namespace, "Service.Name", cassandraService.Name)
			return reconcile.Result{}, err
		}

		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "CASSANDRA Failed to get Deployment")
		return reconcile.Result{}, err
	}

	cnDeployment := r.deploymentForCnV2(instance)
	hssDeployment := r.deploymentForHssV2(instance)
	mmeDeployment := r.deploymentForMmeV2(instance)
	spgwcDeployment := r.deploymentForSpgwcV2(instance)
	spgwuDeployment := r.deploymentForSpgwuV2(instance)
	flexranDeployment := r.deploymentForFlexRANV2(instance)

	cn := &appsv1.Deployment{}
	hss := &appsv1.Deployment{}
	mme := &appsv1.Deployment{}
	spgwc := &appsv1.Deployment{}
	spgwu := &appsv1.Deployment{}
	flexran := &appsv1.Deployment{}
	cnService := r.genCNV2Service(instance)
	hssService := r.genHssV2Service(instance)
	spgwcV2Service := r.genSpgwcV2Service(instance)
	spgwuService := r.genSpgwuV2Service(instance)
	mmeService := r.genMmeV1Service(instance)
	flexranService := r.genFlexranService(instance)
	ranService := r.genRanService(instance)
	if instance.Spec.CoreNetworkAllInOne == true {
		// Creat an oaicn deployment
		// cn := &appsv1.Deployment{}
		// cnDeployment := r.deploymentForCnV2(instance)
		// Check if the oai-cn deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: cnDeployment.GetName(), Namespace: instance.Namespace}, cn)
		if err != nil && errors.IsNotFound(err) {
			if cassandradb.Status.ReadyReplicas == 0 {
				return reconcile.Result{Requeue: true}, Err.New("No cassandradb POD is ready")
			}
			reqLogger.Info("MME domain name " + instance.Spec.MmeDomainName)
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", cnDeployment.Namespace, "Deployment.Name", cnDeployment.Name)
			err = r.client.Create(context.TODO(), cnDeployment)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", cnDeployment.Namespace, "Deployment.Name", cnDeployment.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			// d, _ := time.ParseDuration("20s")
			d, _ := time.ParseDuration("30s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "CN Failed to get Deployment")
			return reconcile.Result{}, err
		}

		// Create an oaicn service
		service := &v1.Service{}
		// cnService := r.genCNV2Service(instance)
		// Check if the oai-cn service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: cnService.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), cnService)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", cnService.Namespace, "Service.Name", cnService.Name)
				return reconcile.Result{}, err
			}
		}
	} else {
		//time.Sleep(15 * time.Second)
		// Creat an oaihss deployment
		// hss := &appsv1.Deployment{}
		// hssDeployment := r.deploymentForHssV2(instance)
		// Check if the oai-hss deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: hssDeployment.GetName(), Namespace: instance.Namespace}, hss)
		if err != nil && errors.IsNotFound(err) {
			if cassandradb.Status.ReadyReplicas == 0 {
				return reconcile.Result{Requeue: true}, Err.New("No cassandradb POD is ready")
			}
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", hssDeployment.Namespace, "Deployment.Name", hssDeployment.Name)
			err = r.client.Create(context.TODO(), hssDeployment)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", hssDeployment.Namespace, "Deployment.Name", hssDeployment.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "HSS Failed to get Deployment")
			return reconcile.Result{}, err
		}
		/////////////////////////////////////////////////////////////////////////////////////////
		//time.Sleep(15 * time.Second)
		// Create an oaihss service
		service := &v1.Service{}
		// hssService := r.genHssV2Service(instance)
		// Check if the oai-hss service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: hssService.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), hssService)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of HSS", "Service.Namespace", hssService.Namespace, "Service.Name", hssService.Name)
				return reconcile.Result{}, err
			}
		}
		///////////////////////////////////////////////////////////////////////////////////////////
		//time.Sleep(15 * time.Second)
		// Creat an oaispgw deployment
		// spgwc := &appsv1.Deployment{}
		// spgwcDeployment := r.deploymentForSpgwcV2(instance)
		// Check if the oai-mme deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: spgwcDeployment.GetName(), Namespace: instance.Namespace}, spgwc)
		if err != nil && errors.IsNotFound(err) {
			if hss.Status.ReadyReplicas == 0 {
				return reconcile.Result{Requeue: true}, Err.New("No hss POD is ready")
			}
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", spgwcDeployment.Namespace, "Deployment.Name", spgwcDeployment.Name)
			err = r.client.Create(context.TODO(), spgwcDeployment)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", spgwcDeployment.Namespace, "Deployment.Name", spgwcDeployment.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "SPGWC Failed to get Deployment")
			return reconcile.Result{}, err
		}
		/////////////////////////////////////////////////////////////////////////////////////////
		//time.Sleep(5 * time.Second)
		// Create an oaispgw service
		// spgwcV2Service := r.genSpgwcV2Service(instance)
		// Check if the oai-spgwc service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: spgwcV2Service.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), spgwcV2Service)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of SPGWC", "Service.Namespace", spgwcV2Service.Namespace, "Service.Name", spgwcV2Service.Name)
				return reconcile.Result{}, err
			}
		}
		///////////////////////////////////////////////////////////////////////////////////////////
		time.Sleep(15 * time.Second)
		// Creat an oaimme deployment
		// mme := &appsv1.Deployment{}
		// mmeDeployment := r.deploymentForMmeV2(instance)
		// Check if the oai-mme deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: mmeDeployment.GetName(), Namespace: instance.Namespace}, mme)
		if err != nil && errors.IsNotFound(err) {
			// if spgwc.Status.ReadyReplicas == 0 {
			// 	return reconcile.Result{Requeue: true}, Err.New("No spgwc POD is ready")
			// }
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", mmeDeployment.Namespace, "Deployment.Name", mmeDeployment.Name)
			err = r.client.Create(context.TODO(), mmeDeployment)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", mmeDeployment.Namespace, "Deployment.Name", mmeDeployment.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "MME Failed to get Deployment")
			return reconcile.Result{}, err
		}
		/////////////////////////////////////////////////////////////////////////////////////////
		//time.Sleep(5 * time.Second)
		// Create an oaimme service
		// mmeService := r.genMmeV1Service(instance)
		// Check if the oai-mme service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: mmeService.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), mmeService)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of MME", "Service.Namespace", mmeService.Namespace, "Service.Name", mmeService.Name)
				return reconcile.Result{}, err
			}
		}
		////////////////////////////////////////////////////////////////////////////////////////////
		//********************************************************************************************************
		time.Sleep(15 * time.Second)
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: spgwuDeployment.GetName(), Namespace: instance.Namespace}, spgwu)
		if err != nil && errors.IsNotFound(err) {
			// if spgwc.Status.ReadyReplicas == 0 {
			// 	return reconcile.Result{Requeue: true}, Err.New("No spgwc POD is ready")
			// }
			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", spgwuDeployment.Namespace, "Deployment.Name", spgwuDeployment.Name)
			err = r.client.Create(context.TODO(), spgwuDeployment)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", spgwuDeployment.Namespace, "Deployment.Name", spgwuDeployment.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully. Let's wait for it to be ready
			d, _ := time.ParseDuration("5s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		} else if err != nil {
			reqLogger.Error(err, "SPGWU Failed to get Deployment")
			return reconcile.Result{}, err
		}

		// Check if the oai-spgwu service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: spgwuService.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), spgwuService)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service of SPGWU", "Service.Namespace", spgwuService.Namespace, "Service.Name", spgwuService.Name)
				return reconcile.Result{}, err
			}
		}
		//********************************************************************************************************
	}
	/////////////////////////////////////////////////////////////////////////////
	if instance.Spec.FlexRAN == true {
		// Check if the flexran deployment already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: flexranDeployment.GetName(), Namespace: instance.Namespace}, flexran)
		if err != nil && errors.IsNotFound(err) {

			// if mme.Status.ReadyReplicas == 0 {
			// 	d, _ := time.ParseDuration("10s")
			// 	return reconcile.Result{Requeue: true, RequeueAfter: d}, Err.New("No oai-mme POD is ready, 10 seconds backoff")
			// }

			reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", flexranDeployment.Namespace, "Deployment.Name", flexranDeployment.Name)
			err = r.client.Create(context.TODO(), flexranDeployment)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", flexranDeployment.Namespace, "Deployment.Name", flexranDeployment.Name)
				return reconcile.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "FlexRAN Failed to get Deployment")
			return reconcile.Result{}, err
		}

		// Create an FlexRAN service
		service := &v1.Service{}
		// Check if the flexran service already exists, if not create a new one
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: flexranService.GetName(), Namespace: instance.Namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			err = r.client.Create(context.TODO(), flexranService)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", flexranService.Namespace, "Service.Name", flexranService.Name)
				return reconcile.Result{}, err
			}
		}

		time.Sleep(5 * time.Second)
	}
	/////////////////////////////////////////////////////////////////////////////////
	// time.Sleep(15 * time.Second)
	// Create an oairan deployment
	ran := &appsv1.Deployment{}
	ranDeployment := r.deploymentForRAN(instance)
	// Check if the oai-ran deployment already exists, if not create a new one
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ranDeployment.GetName(), Namespace: instance.Namespace}, ran)
	if err != nil && errors.IsNotFound(err) {
		// if mme.Status.ReadyReplicas == 0 {
		// 	d, _ := time.ParseDuration("10s")
		// 	return reconcile.Result{Requeue: true, RequeueAfter: d}, Err.New("No oai-mme POD is ready, 10 seconds backoff")
		// }
		if instance.Spec.CoreNetworkAllInOne == true {
			if cn.Status.ReadyReplicas == 0 {
				// d, _ := time.ParseDuration("5s")
				d, _ := time.ParseDuration("10s")
				return reconcile.Result{Requeue: true, RequeueAfter: d}, Err.New("No oai-cn POD is ready, 10 seconds backoff")
			}
		} else {
			if mme.Status.ReadyReplicas == 0 {
				// d, _ := time.ParseDuration("5s")
				d, _ := time.ParseDuration("10s")
				return reconcile.Result{Requeue: true, RequeueAfter: d}, Err.New("No oai-mme POD is ready, 10 seconds backoff")
			}
		}
		reqLogger.Info("Sheeps are ready")
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
	////////////////////////////////////////////
	// Create an oairan service
	service := &v1.Service{}
	// ranService := r.genRanService(instance)
	// Check if the oai-cn service already exists, if not create a new one
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ranService.GetName(), Namespace: instance.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		err = r.client.Create(context.TODO(), ranService)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", ranService.Namespace, "Service.Name", ranService.Name)
			return reconcile.Result{}, err
		}
	}

	// Ensure the deployment size is the same as the spec
	// size := instance.Spec.Size
	if instance.Spec.CoreNetworkAllInOne == true {
		size := instance.Spec.OaiCnSize
		if *cn.Spec.Replicas != size {
			cn.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), cn)
			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", cn.Namespace, "Deployment.Name", cn.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	} else {
		// Ensure the deployment size is the same as the spec
		size := instance.Spec.OaiHssSize
		if *hss.Spec.Replicas != size {
			hss.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), hss)
			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", hss.Namespace, "Deployment.Name", hss.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
		// mme
		size = instance.Spec.OaiMmeSize
		if *mme.Spec.Replicas != size {
			mme.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), mme)
			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", mme.Namespace, "Deployment.Name", mme.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
		// spgwc
		size = instance.Spec.OaiSpgwcSize
		if *spgwc.Spec.Replicas != size {
			spgwc.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), spgwc)
			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", spgwc.Namespace, "Deployment.Name", spgwc.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}

		// spgwu
		size = instance.Spec.OaiSpgwuSize
		if *spgwu.Spec.Replicas != size {
			spgwu.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), spgwu)
			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", spgwu.Namespace, "Deployment.Name", spgwu.Name)
				return reconcile.Result{}, err
			}
			// Spec updated - return and requeue
			return reconcile.Result{Requeue: true}, nil
		}
	}
	if instance.Spec.FlexRAN == true {
		// flexran
		size := instance.Spec.FlexranSize
		if *flexran.Spec.Replicas != size {
			flexran.Spec.Replicas = &size
			err = r.client.Update(context.TODO(), flexran)
			if err != nil {
				reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", flexran.Namespace, "Deployment.Name", flexran.Name)
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

	// Check configmap is fine or not. If it's changed, update ConfigMap and restart cn ran
	if err == nil {
		if reflect.DeepEqual(new.Data, config.Data) {
			reqLogger.Info("newConf equals config")

		} else {
			reqLogger.Info("newConf does not equals config")
			reqLogger.Info("Update ConfigMap and deleting CN and RAN")

			//****************************************
			// // fine grained upate; update only the concerned entities
			// Mosaic5gSpec := mosaic5gv1alpha1.Mosaic5gSpec{}
			// if err := yaml.Unmarshal([]byte(new.Data["conf.yaml"]), &Mosaic5gSpec); err != nil {
			// 	panic(err)
			// } else {

			// 	fmt.Println("success convert string to yaml")
			// 	fmt.Println(reflect.TypeOf(Mosaic5gSpec))
			// 	fmt.Println(Mosaic5gSpec.CoreNetworkAllInOne)

			// }
			//****************************************
			err = r.client.Update(context.TODO(), new)
			//Should only kill the POD
			/////////////////////////////////////////////////////////////////
			// err = r.client.Delete(context.TODO(), cnDeployment)
			// err = r.client.Delete(context.TODO(), cnService)
			// err = r.client.Delete(context.TODO(), hssDeployment)
			// err = r.client.Delete(context.TODO(), hssService)
			// err = r.client.Delete(context.TODO(), mmeDeployment)
			// err = r.client.Delete(context.TODO(), mmeService)
			// err = r.client.Delete(context.TODO(), spgwcDeployment)
			// err = r.client.Delete(context.TODO(), spgwcV2Service)
			err = r.client.Delete(context.TODO(), ranDeployment)
			err = r.client.Delete(context.TODO(), ranService)
			/////////////////////////////////////////////////////////////////
			// var sizeReset int32 = 0
			// if instance.Spec.CoreNetworkAllInOne == true {
			// 	// cn
			// 	if *cn.Spec.Replicas != sizeReset {
			// 		cn.Spec.Replicas = &sizeReset
			// 		err = r.client.Update(context.TODO(), cn)
			// 		if err != nil {
			// 			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", cn.Namespace, "Deployment.Name", cn.Name)
			// 			return reconcile.Result{}, err
			// 		}
			// 		// Spec updated - return and requeue
			// 		return reconcile.Result{Requeue: true}, nil
			// 	}
			// } else {
			// 	// hss
			// 	if *hss.Spec.Replicas != sizeReset {
			// 		hss.Spec.Replicas = &sizeReset
			// 		err = r.client.Update(context.TODO(), hss)
			// 		if err != nil {
			// 			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", hss.Namespace, "Deployment.Name", hss.Name)
			// 			return reconcile.Result{}, err
			// 		}
			// 		// Spec updated - return and requeue
			// 		return reconcile.Result{Requeue: true}, nil
			// 	}
			// 	// mme
			// 	if *mme.Spec.Replicas != sizeReset {
			// 		mme.Spec.Replicas = &sizeReset
			// 		err = r.client.Update(context.TODO(), mme)
			// 		if err != nil {
			// 			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", mme.Namespace, "Deployment.Name", mme.Name)
			// 			return reconcile.Result{}, err
			// 		}
			// 		// Spec updated - return and requeue
			// 		return reconcile.Result{Requeue: true}, nil
			// 	}
			// 	// spgwc
			// 	if *spgwc.Spec.Replicas != sizeReset {
			// 		spgwc.Spec.Replicas = &sizeReset
			// 		err = r.client.Update(context.TODO(), spgwc)
			// 		if err != nil {
			// 			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", spgwc.Namespace, "Deployment.Name", spgwc.Name)
			// 			return reconcile.Result{}, err
			// 		}
			// 		// Spec updated - return and requeue
			// 		return reconcile.Result{Requeue: true}, nil
			// 	}
			// }
			/////////////////////////////////////////////////////////////////
			// Spec updated - return and requeue
			// d, _ := time.ParseDuration("5s")
			d, _ := time.ParseDuration("10s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		}

	}
	// Everything is fine, Reconcile ends
	return reconcile.Result{}, nil
}

// deploymentForHssV2 returns a HSS Network Deployment object
func (r *ReconcileMosaic5g) deploymentForHssV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {

	hssName := m.Spec.HssDomainName
	//ls := util.LabelsForMosaic5g(m.Name + hssName)
	replicas := m.Spec.OaiHssSize
	labels := make(map[string]string)
	labels["app"] = "oaihss"
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/oaihss"] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.GetName() + "-" + hssName,
			Namespace:   m.Namespace,
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
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiHssImage,
						Name:            "oaihss",
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
						},
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
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForMmeV2 returns a MME Network Deployment object
func (r *ReconcileMosaic5g) deploymentForMmeV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {

	mmeName := m.Spec.MmeDomainName
	//ls := util.LabelsForMosaic5g(m.Name + mmeName)
	replicas := m.Spec.OaiMmeSize
	labels := make(map[string]string)
	labels["app"] = "oaimme"
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/oaimme"] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.GetName() + "-" + mmeName,
			Namespace:   m.Namespace,
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
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiMmeImage,
						Name:            "oaimme",
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
						},
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
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForSpgwcV2 returns a SPGWC Network Deployment object
func (r *ReconcileMosaic5g) deploymentForSpgwcV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {

	spgwcName := m.Spec.SpgwcDomainName
	//ls := util.LabelsForMosaic5g(m.Name + spgwcName)
	replicas := m.Spec.OaiSpgwcSize
	labels := make(map[string]string)
	labels["app"] = "oaispgwc"
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/oaispgwc"] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.GetName() + "-" + spgwcName,
			Namespace:   m.Namespace,
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
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiSpgwcImage,
						Name:            "oaispgwc",
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
						},
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
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForSpgwuV2 returns a SPGWU Network Deployment object
func (r *ReconcileMosaic5g) deploymentForSpgwuV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {

	spgwuName := m.Spec.SpgwuDomainName
	//ls := util.LabelsForMosaic5g(m.Name + spgwuName)
	replicas := m.Spec.OaiSpgwuSize
	labels := make(map[string]string)
	labels["app"] = "oaispgwu"
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/oaispgwu"] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.GetName() + "-" + spgwuName,
			Namespace:   m.Namespace,
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
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.OaiSpgwuImage,
						Name:            "oaispgwu",
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
						},
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
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForCnV2 returns a Core Network Deployment object
func (r *ReconcileMosaic5g) deploymentForCnV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {

	cnName := m.Spec.MmeDomainName
	//ls := util.LabelsForMosaic5g(m.Name + cnName)
	replicas := m.Spec.OaiCnSize
	labels := make(map[string]string)
	labels["app"] = "oaicn"
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/oaicn"] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.GetName() + "-" + cnName,
			Namespace:   m.Namespace,
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
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Containers: []corev1.Container{{
						Image:           m.Spec.CNImage,
						Name:            "oaicn",
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
						},
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
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForRAN returns a Core Network Deployment object
func (r *ReconcileMosaic5g) deploymentForRAN(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	// ls := util.LabelsForMosaic5g(m.Name)
	replicas := m.Spec.OaiRanSize
	labels := make(map[string]string)
	labels["app"] = "oaienb"
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+m.Name+"-"+"oairan"] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.GetName() + "-" + "oairan",
			Namespace:   m.Namespace,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				// MatchLabels: ls,
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					// Labels: ls,
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					// NodeSelector: map[string]string{
					// 	"usrp": "true"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.RANImage,
						Name:            "oairan",
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
						},
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
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForCassandra returns a Cassandra Deployment object
func (r *ReconcileMosaic5g) deploymentForCassandra(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	//ls := util.LabelsForMosaic5g(m.Name + cnName)
	var replicas int32
	replicas = m.Spec.CassandraSize
	labels := make(map[string]string)
	labels["app"] = "cassandra"
	// labels["app"] = m.Spec.CassandraDomainName
	// Annotations := make(map[string]string)
	// Annotations["container.apparmor.security.beta.kubernetes.io/cassandra"] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.GetName() + "-" + m.Spec.CassandraDomainName,
			Namespace: m.Namespace,
			Labels:    labels,
			// Annotations: Annotations,
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
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image: m.Spec.CassandraImage,
						Name:  "cassandra",
						// Env: []corev1.EnvVar{
						// 	{Name: "MYSQL_ROOT_PASSWORD", Value: "linux"},
						// },
						Ports: []corev1.ContainerPort{{
							Name:          "cassandra-1",
							ContainerPort: 7000,
							Protocol:      corev1.ProtocolTCP,
						}, {
							Name:          "cassandra-2",
							ContainerPort: 7001,
							Protocol:      corev1.ProtocolTCP,
						}, {
							Name:          "cassandra-3",
							ContainerPort: 7199,
							Protocol:      corev1.ProtocolTCP,
						}, {
							Name:          "cassandra-4",
							ContainerPort: 9042,
							Protocol:      corev1.ProtocolTCP,
						}, {
							Name:          "cassandra-5",
							ContainerPort: 9160,
							Protocol:      corev1.ProtocolTCP,
						}},
					}},
					Affinity: util.GenAffinity("cn"),
				},
			},
		},
	}
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForFlexRANV2 returns a FlexRAN Deployment object
func (r *ReconcileMosaic5g) deploymentForFlexRANV2(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	var replicas int32
	replicas = m.Spec.FlexranSize
	labels := make(map[string]string)
	// labels["app"] = m.Spec.FlexRANDomainName
	labels["app"] = "flexran"
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/flexran"] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.GetName() + "-" + m.Spec.FlexRANDomainName,
			Namespace:   m.Namespace,
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
					// NodeSelector: map[string]string{
					// 	"usrp": "false"},
					Hostname: "ubuntu",
					Containers: []corev1.Container{{
						Image:           m.Spec.FlexranImage,
						Name:            "flexran",
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("10m"),
								corev1.ResourceMemory: resource.MustParse("250Mi"),
							},
						},
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
						// Env: []corev1.EnvVar{
						// 	{Name: "MYSQL_ROOT_PASSWORD", Value: "linux"},
						// },
						Ports: []corev1.ContainerPort{{
							Name:          "sbi-port",
							ContainerPort: 2210,
							// Protocol:      corev1.ProtocolTCP,
						}, {
							Name:          "nbi-port",
							ContainerPort: 9999,
							// Protocol:      corev1.ProtocolTCP,
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

// genCNV2Service will generate a service for oaicn
func (r *ReconcileMosaic5g) genCNV2Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	var service *v1.Service
	selectMap := make(map[string]string)
	selectMap["app"] = "oaicn"
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
	service.Name = "oaicn"
	service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genHssV2Service will generate a service for oaicn
func (r *ReconcileMosaic5g) genHssV2Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	var service *v1.Service
	selectMap := make(map[string]string)
	selectMap["app"] = "oaihss"
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
	service.Name = "oaihss"
	service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genMmeV1Service will generate a service for oaicn
func (r *ReconcileMosaic5g) genMmeV1Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	var service *v1.Service
	selectMap := make(map[string]string)
	selectMap["app"] = "oaimme"
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
	service.Name = "oaimme"
	service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genSpgwcV2Service will generate a service for oaicn
func (r *ReconcileMosaic5g) genSpgwcV2Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	var service *v1.Service
	selectMap := make(map[string]string)
	selectMap["app"] = "oaispgwc"
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
	service.Name = "oaispgwc"
	service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genSpgwuV2Service will generate a service for oaicn
func (r *ReconcileMosaic5g) genSpgwuV2Service(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	var service *v1.Service
	selectMap := make(map[string]string)
	selectMap["app"] = "oaispgwu"
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
		Type:     "NodePort",
		// ClusterIP: "None",
	}
	service.Name = "oaispgwu"
	service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genRanService will generate a service for oaicn
func (r *ReconcileMosaic5g) genRanService(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	var service *v1.Service
	selectMap := make(map[string]string)
	selectMap["app"] = "oairan"
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
	service.Name = "oairan"
	service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genCassandraService will generate a service for oaicn
func (r *ReconcileMosaic5g) genCassandraService(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	var service *v1.Service
	selectMap := make(map[string]string)
	selectMap["app"] = m.Spec.CassandraDomainName
	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{Name: "cassandra-port-1",
				Protocol: corev1.ProtocolTCP,
				Port:     7000,
			}, {Name: "cassandra-port-2",
				Protocol: corev1.ProtocolTCP,
				Port:     7001,
			}, {
				Name:     "cassandra-3",
				Port:     7199,
				Protocol: corev1.ProtocolTCP,
			}, {
				Name:     "cassandra-4",
				Port:     9042,
				Protocol: corev1.ProtocolTCP,
			}, {
				Name:     "cassandra-5",
				Port:     9160,
				Protocol: corev1.ProtocolTCP,
			},
		},
		Selector: selectMap,
		Type:     "NodePort",
		// ClusterIP: "None",
	}
	service.Name = m.Spec.CassandraDomainName
	service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genFlexranService
// genFlexranService will generate a service for flexran
func (r *ReconcileMosaic5g) genFlexranService(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	var service *v1.Service
	selectMap := make(map[string]string)
	selectMap["app"] = "flexran"
	// selectMap["app"] = m.Spec.FlexRANDomainName
	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{
				Name: "sbi-port",
				Port: 2210,
				// Protocol:      corev1.ProtocolTCP,
			}, {
				Name: "nbi-port",
				Port: 9999,
				// Protocol:      corev1.ProtocolTCP,
			},
		},
		Selector: selectMap,
		Type:     "NodePort",
		// ClusterIP: "None",
	}
	service.Name = "flexran"
	// service.Name = m.Spec.FlexRANDomainName
	service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}
