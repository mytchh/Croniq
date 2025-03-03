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