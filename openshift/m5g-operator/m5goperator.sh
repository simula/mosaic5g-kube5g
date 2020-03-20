#!/bin/bash

# prepare ENVs
#export KUBECONFIG=/home/agrion/kubernetes/aiyu
export OPERATOR_NAME=m5g-operator
export MYNAME=${USER}
export MYDNS="192.168.1.1"


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
    operator-sdk up local --namespace=default
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
    esac
}

deploy_operator_from_clean_machine(){
    echo "Start a fresh microk8s and deploy operator on it, tested with Ubuntu 18.04"
    echo "sudo without password is recommended"
    sudo snap install microk8s --classic --channel=1.14/stable
    sudo snap install kubectl --classic
    microk8s.start
    microk8s.enable dns
    # kubeconfig is used by operator
    sudo chown ${MYNAME} -R $HOME/.kube
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
    sleep 3
    echo "Done, now run [local] or [container start] to create your operator"
}

clean(){
    kubectl delete -f deploy/crds/mosaic5g_v1alpha1_mosaic5g_crd.yaml
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
        -i | --install)
            deploy_operator_from_clean_machine 
        ;;
        -r | --remove)
            break_down
        ;;
        *)
            echo_info '
This program installs the requirements to run kubernets on one machine, 
and run mosaic5g-operator as a pod inside kubernetes in order to manage
the deployments and services of 4G/5G networks in cloud native environment.
This program also allows to run mosaic5g-operator locally as Golang app.
Options:
-i | --install
Install and run microk8s kubectl, then deploy operator on it"
-n | --init
Apply CRD to k8s cluster (Required for Operator)"
-l | --local
Run Operator as a Golang app at local"
container [start|stop]
Run Operator as a POD inside Kubernetes"
-c | --clean 
Remove CRD from cluster"
-r | --remove
remove the snap of kubectl and microk8s
Usage:
./m5goperator.sh -i 
./m5goperator.sh container start
            '
        ;;
    esac

}
main "$@"
