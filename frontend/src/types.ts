export interface CronJob {
    metadata: {
        name: string;
        namespace: string;
    };
    spec: {
        schedule: string;
        suspend?: boolean;
    };
}

export interface CronJobResponse {
    items: CronJob[];
}

export interface ClusterInfo {
    name: string;
    serverAddress: string;
    version: string;
}

export interface Job {
    metadata: {
        name: string;
        namespace: string;
        creationTimestamp: string;
    };
    status: {
        active: number;
        succeeded: number;
        failed: number;
        startTime?: string;
        completionTime?: string;
    };
}

export interface JobResponse {
    items: Job[];
}

export interface JobStats {
    totalCronJobs: number;
    activeCronJobs: number;
    totalJobs: number;
    runningJobs: number;
    failedJobs: number;
    succeededJobs: number;
} 