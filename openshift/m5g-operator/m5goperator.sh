#!/bin/bash

# prepare ENVs
#export KUBECONFIG=/home/agrion/kubernetes/aiyu
export OPERATOR_NAME=m5g-operator

DOCKER_OPERATOR_REP_NAME="mosaic5gecosys/m5g_operator"
DOCKER_OPERATOR_TAG="v1.test"

###################################
# colorful echos
###################################

black='\E[30m'
red='\E[31m'
green='\E[32m'
yellow='\E[33m'
blue='\E[1;34m'
magenta='\E[35m'
cyan='\E[36m'
white='\E[37m'
reset_color='\E[00m'
COLORIZE=1

cecho()  {  
    # Color-echo
    # arg1 = message
    # arg2 = color
    local default_msg="No Message."
    message=${1:-$default_msg}
    color=${2:-$green}
    [ "$COLORIZE" = "1" ] && message="$color$message$reset_color"
    echo -e "$message"
    return
}

echo_error()   { cecho "$*" $red          ;}
echo_fatal()   { cecho "$*" $red; exit -1 ;}
echo_warn()    { cecho "$*" $yellow       ;}
echo_success() { cecho "$*" $green        ;}
echo_info()    { cecho "$*" $blue         ;}

run_local(){
    # operator-sdk run --local --namespace=default
    operator-sdk run local --watch-namespace=default
}

run_container(){
    case ${1} in
        start)
            kubectl apply -f deploy/service_account.yaml
            kubectl apply -f deploy/role.yaml
            kubectl apply -f deploy/role_binding.yaml
            kubectl apply -f deploy/operator.yaml
        ;;
        stop)
            kubectl delete -f deploy/service_account.yaml
            kubectl delete -f deploy/role.yaml
            kubectl delete -f deploy/role_binding.yaml
            kubectl delete -f deploy/operator.yaml
        ;;
        *)
            echo_error "Unkown option '${1}' for container"
    esac
}

operator_gen_build(){
    case ${1} in
        -g | --generate)
            echo_info "generating the crds of the operator"
            operator-sdk generate k8s
            operator-sdk generate crds
        ;;
        #######################
        -b | --build)
            docker_operator_repository_name=""
            docker_operator_tag=""
            case ${2} in
                -d | --default)
                    docker_operator_repository_name=$DOCKER_OPERATOR_REP_NAME
                ;;
                *)
                    docker_operator_repository_name=${2}
            esac
            case ${3} in
                -d | --default)
                    docker_operator_tag=$DOCKER_OPERATOR_TAG
                ;;
                *)
                    docker_operator_tag=${3}
            esac
            docker_operator_image_tag=${docker_operator_repository_name}":"${docker_operator_tag}
            echo_info "building the docker image $docker_operator_image_tag for the operator"
            operator-sdk build $docker_operator_image_tag
            # operator-sdk build mosaic5gecosys/m5g_operator:1.1
            echo_success "the docker image of the operator $docker_operator_image_tag is successfully build"
            echo_info "do not forgot to push it to docker hub by: docker push $docker_operator_image_tag"
        #######################
        ;;
        -p | --push)
            docker_operator_repository_name=""
            docker_operator_tag=""
            case ${2} in
                -d | --default)
                    docker_operator_repository_name=$DOCKER_OPERATOR_REP_NAME
                ;;
                *)
                    docker_operator_repository_name=${2}
            esac
            case ${3} in
                -d | --default)
                    docker_operator_tag=$DOCKER_OPERATOR_TAG
                ;;
                *)
                    docker_operator_tag=${3}
            esac
            docker_operator_image_tag=${docker_operator_repository_name}":"${docker_operator_tag}
            echo_info "pushing the docker image $docker_operator_image_tag to docker hub"
            docker push $docker_operator_image_tag
            echo_success "the docker image of the operator $docker_operator_image_tag is successfully pushed to docker hub"
        #######################
        ;;
        *)
            echo_error "Unkown option '${1}' for operator operations"
    esac
}
apply_cr(){
    case ${2} in
        aio|all-in-one)
            echo_info "kubectl apply -f deploy/crds/cr-${1}/lte-all-in-one/mosaic5g_v1alpha1_cr_${1}_lte_all_in_one.yaml"
            kubectl apply -f deploy/crds/cr-${1}/lte-all-in-one/mosaic5g_v1alpha1_cr_${1}_lte_all_in_one.yaml
            echo "lte network Custom Resources (CR) of monolitic Core Network, of snap version ${1}, is applied"
        ;;
        dis-cn|disaggregated-cn)
            echo_info "kubectl apply -f deploy/crds/cr-${1}/lte/mosaic5g_v1alpha1_cr_${1}_lte.yaml"
            kubectl apply -f deploy/crds/cr-${1}/lte/mosaic5g_v1alpha1_cr_${1}_lte.yaml
            echo "lte network Custom Resources (CR) of disaggregated Core Network entities, of snap version ${1}, is applied"
        ;;
        *)
            echo_error "Unkown option '${1}' for deploy"
    esac
}

downgrade_image(){
    APISERVER=`kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}'`
    TOKEN=`kubectl get secret $(kubectl get serviceaccount default -o jsonpath='{.secrets[0].name}') -o jsonpath='{.data.token}' | base64 --decode `
    
    curl \
    -H "content-Type: application/json-patch+json" \
    -H "Authorization: Bearer ${TOKEN}"\
    --insecure \
    -X PATCH ${APISERVER}/apis/mosaic5g.com/v1alpha1/namespaces/default/mosaic5gs/mosaic5g \
    -d '[{"op":"replace","path":"/spec/cnImage","value":"arouk/oaicn:1.0"},{"op":"replace","path":"/spec/ranImage","value":"arouk/oairan:1.0"}]'
    echo " "
    echo "Core Network is downgraded to version 1.0, and RAN network is downgraded to 1.0"
}

upgrade_image(){
    APISERVER=`kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}'`
    TOKEN=`kubectl get secret $(kubectl get serviceaccount default -o jsonpath='{.secrets[0].name}') -o jsonpath='{.data.token}' | base64 --decode `

    curl \
    -H "content-Type: application/json-patch+json" \
    -H "Authorization: Bearer ${TOKEN}"\
    --insecure \
    -X PATCH ${APISERVER}/apis/mosaic5g.com/v1alpha1/namespaces/default/mosaic5gs/mosaic5g \
    -d '[{"op":"replace","path":"/spec/cnImage","value":"arouk/oaicn:1.1"},{"op":"replace","path":"/spec/ranImage","value":"arouk/oairan:1.1"}]'
    echo " "
    echo "Core Network is upgraded to version 1.1, and RAN network is upgraded to 1.1"
}

delete_cr(){
    kubectl delete -f deploy/crds/mosaic5g_v1alpha1_mosaic5g_cr.yaml
    echo "Custom Resources (CR) of the network is deleted"
}

deploy_operator_from_clean_machine(){
    echo "Start a fresh microk8s and deploy operator on it, tested with Ubuntu 18.04"
    echo "sudo without password is recommended"
    sudo snap install microk8s --classic --channel=1.14/stable
    sudo snap install kubectl --classic
    microk8s.start
    microk8s.enable dns
    # kubeconfig is used by operator
    sudo chown $USER -R $HOME/.kube
    microk8s.kubectl config view --raw > $HOME/.kube/config
    # enable privileged
    sudo bash -c 'echo "--allow-privileged=true" >> /var/snap/microk8s/current/args/kubelet'
    sudo bash -c 'echo "--allow-privileged=true" >> /var/snap/microk8s/current/args/kube-apiserver'
    # Restart kube
    sudo systemctl restart snap.microk8s.daemon-kubelet.service
    sudo systemctl restart snap.microk8s.daemon-apiserver.service
    # Configure DNS if it's not working 
    # microk8s.kubectl -n kube-system edit configmap/coredns

}

init(){
    echo "Applying crd..."
    kubectl apply -f deploy/crds/mosaic5g_v1alpha1_mosaic5g_crd.yaml
    kubectl apply -f defaultRole.yaml
    sleep 3
    echo "Done, now run [local] or [container start] to create your operator"
}

clean(){
    kubectl delete -f deploy/crds/mosaic5g_v1alpha1_mosaic5g_crd.yaml
    kubectl delete -f defaultRole.yaml
}

break_down(){
    sudo snap remove microk8s 
    sudo snap remove kubectl 
}

main() {
    case ${1} in
        -n | --init)
            init
        ;;
        -c | --clean)
            clean
        ;;
        -l | --local)
            run_local
        ;;
        container)
            run_container ${2}
        ;;
        deploy)
            apply_cr ${2} ${3}
        ;;
        upgrade)
            upgrade_image
        ;;
        downgrade)
            downgrade_image
        ;;
        -d | --delete)
            delete_cr 
        ;;
        -i | --install)
            deploy_operator_from_clean_machine 
        ;;
        -r | --remove)
            break_down
        ;;
        ###
        -o | --operator)
            operator_gen_build ${2} ${3} ${4}
        ;;
        ###
        *)
            echo_info '
This program installs the requirements to run kubernets on one machine, 
and run mosaic5g-operator as a pod inside kubernetes in order to manage
the deployments and services of 4G/5G networks in cloud native environment.
This program also allows to run mosaic5g-operator locally as Golang app.
Options:
-i | --install
    Install and run microk8s kubectl, then deploy operator on it
-n | --init
    Apply CRD to k8s cluster (Required for Operator)
-l | --local
    Run Operator as a Golang app at local
container [start|stop]
    Run Operator as a POD inside Kubernetes
deploy [v1|v2][[aio|all-in-one]|[dis-cn|disaggregated-cn]]
    Deploy the network with:
        v1: snap version v1
        v2: snap version v2
        - aio|all-in-one: all the core network entities (oai-hss, oai-mme, oai-spgw) in one pod
        - dis-cn|disaggregated-cn: the core network entities (oai-hss, oai-mme, oai-spgw) are deployed on disaggregated pods
-d | --delete 
    Stop the network by deleting the Custom Resources (CR) of the network
-c | --clean 
    Remove CRD from cluster
-r | --remove
    remove the snap of kubectl and microk8s
-o | --operator [-g|-b -d -d]
    -g|--generate: generate the crds of the operator
    -b|--build: build docker image of the operator
        -b -d -d: 
                with default values of (docker-hub accout)/(docker-image-name):mosaic5gecosys/m5g_operator, for the first (-d)
                with default tag (1.1) for the second (-d)
            Example: the options "-b -d -d" will build the image "mosaic5gecosys/m5g_operator:1.1"
Usage:
    ./m5goperator.sh -i 
    ./m5goperator.sh -o -g # generate the crds of the operator
    ./m5goperator.sh -o -b -d -d # build the docker image of the operator with the default values: mosaic5gecosys/m5g_operator:1.1
    ./m5goperator.sh -o -b -d v1.test # build the docker image of the operator with the default values of docker image name and certain tage: mosaic5gecosys/m5g_operator:v1.test

    ./m5goperator.sh container start # start kube5g-operator as pod in K8S
    ./m5goperator.sh deploy v1 all-in-one # deploy network with all-in-one for the core, for the snap version v1
    ./m5goperator.sh deploy v2 all-in-one # deploy network with all-in-one for the core, for the snap version v2
    ./m5goperator.sh deploy v1 disaggregated-cn # deploy network with disaggregated core, for the snap version v1
    ./m5goperator.sh deploy v2 disaggregated-cn # deploy network with disaggregated core, for the snap version v2
            '
        ;;
    esac

}
main "$@"
