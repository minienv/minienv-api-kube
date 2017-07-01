package main

import (
	"fmt"
	"log"
	"strings"
)

var VAR_JOB_NAME string = "$jobName"

func isProvisionerRunning(envId string, kubeServiceToken string, kubeServiceBaseUrl string, kubeNamespace string) (bool, error) {
	name, err := getPodName(getProvisionerJobName(envId), kubeServiceToken, kubeServiceBaseUrl, kubeNamespace)
	if err != nil {
		log.Println("Error checking provisioner status: ", err)
		return false, err
	} else {
		return name != "", nil
	}
}

func deleteProvisioner(envId string, kubeServiceToken string, kubeServiceBaseUrl string, kubeNamespace string) (bool, error) {
	return deleteJob(getProvisionerJobName(envId), kubeServiceToken, kubeServiceBaseUrl, kubeNamespace)
}

func deployProvisioner(envId string, storageDriver string, pvTemplate string, pvcTemplate string, jobTemplate string, kubeServiceToken string, kubeServiceBaseUrl string, kubeNamespace string) (error) {
	// delete example, if it exists
	deleteProvisioner(envId, kubeServiceToken, kubeServiceBaseUrl, kubeNamespace)
	// create persistent volume if not exits
	pvName := getPersistentVolumeName(envId)
	pvPath := getPersistentVolumePath(envId)
	pvResponse, err := getPersistentVolume(pvName, kubeServiceToken, kubeServiceBaseUrl)
	if err != nil {
		log.Println("Error saving persistent volume: ", err)
		return err
	} else if pvResponse == nil {
		pv := pvTemplate
		pv = strings.Replace(pv, VAR_PV_NAME, pvName, -1)
		pv = strings.Replace(pv, VAR_PV_PATH, pvPath, -1)
		_, err = savePersistentVolume(pv, kubeServiceToken, kubeServiceBaseUrl)
		if err != nil {
			log.Println("Error saving persistent volume: ", err)
			return err
		}
	}
	// create persistent volume claim, if not exists
	pvcName := getPersistentVolumeClaimName(envId)
	pvcResponse, err := getPersistentVolumeClaim(pvcName, kubeServiceToken, kubeServiceBaseUrl, kubeNamespace)
	if err != nil {
		log.Println("Error saving persistent volume claim: ", err)
		return err
	} else if pvcResponse == nil {
		pvc := pvcTemplate
		pvc = strings.Replace(pvc, VAR_PVC_NAME, pvcName, -1)
		_, err = savePersistentVolumeClaim(pvc, kubeServiceToken, kubeServiceBaseUrl, kubeNamespace)
		if err != nil {
			log.Println("Error saving persistent volume claim: ", err)
			return err
		}
	}
	// create job
	jobName := getProvisionerJobName(envId)
	job := jobTemplate
	job = strings.Replace(job, VAR_JOB_NAME, jobName, -1)
	job = strings.Replace(job, VAR_STORAGE_DRIVER, storageDriver, -1)
	job = strings.Replace(job, VAR_PVC_NAME, pvcName, -1)
	_, err = saveJob(job, kubeServiceToken, kubeServiceBaseUrl, kubeNamespace)
	if err != nil {
		log.Println("Error saving job: ", err)
		return err
	}
	return nil
}

func getProvisionerJobName(envId string) string {
	return strings.ToLower(fmt.Sprintf("env-%s-provision-job", envId))
}