import { CronJob, CronJobResponse, ClusterInfo } from './types.ts';

class CronJobUI {
    private tableBody: HTMLElement;
    private errorDiv: HTMLElement;
    private clusterInfoDiv: HTMLElement;
    private debugDiv: HTMLElement;

    constructor() {
        this.tableBody = document.getElementById('cronJobList') as HTMLElement;
        this.errorDiv = document.getElementById('error') as HTMLElement;
        this.clusterInfoDiv = document.getElementById('clusterInfo') as HTMLElement;
        this.debugDiv = document.getElementById('debug') as HTMLElement;
        this.init();
    }

    private async init(): Promise<void> {
        try {
            await this.fetchClusterInfo();
            await this.fetchAndDisplayCronJobs();
        } catch (err) {
            this.showError(err instanceof Error ? err.message : 'An error occurred');
        }
    }

    private async fetchClusterInfo(): Promise<void> {
        try {
            const response = await fetch('http://localhost:8080/api/cluster-info');
            this.log(`Cluster info response status: ${response.status}`);
            
            if (!response.ok) {
                const text = await response.text();
                throw new Error(`Failed to fetch cluster info: ${response.status} ${text}`);
            }
            
            const info: ClusterInfo = await response.json();
            this.log(`Cluster info: ${JSON.stringify(info, null, 2)}`);
            this.renderClusterInfo(info);
        } catch (err) {
            this.showError(`Cluster info error: ${err instanceof Error ? err.message : String(err)}`);
        }
    }

    private renderClusterInfo(info: ClusterInfo): void {
        this.clusterInfoDiv.innerHTML = `
            <div class="cluster-info">
                <p>Cluster: ${this.escapeHtml(info.name)}</p>
                <p>API Server: ${this.escapeHtml(info.serverAddress)}</p>
                <p>Version: ${this.escapeHtml(info.version)}</p>
            </div>
        `;
    }

    private async fetchAndDisplayCronJobs(): Promise<void> {
        try {
            const response = await fetch('http://localhost:8080/api/cronjobs');
            this.log(`CronJobs response status: ${response.status}`);
            
            if (!response.ok) {
                const text = await response.text();
                throw new Error(`Failed to fetch cron jobs: ${response.status} ${text}`);
            }
            
            const data: CronJobResponse = await response.json();
            this.log(`CronJobs data: ${JSON.stringify(data, null, 2)}`);
            this.renderCronJobs(data.items || []);
        } catch (err) {
            this.showError(`CronJobs error: ${err instanceof Error ? err.message : String(err)}`);
        }
    }

    private renderCronJobs(cronJobs: CronJob[]): void {
        this.tableBody.innerHTML = cronJobs.map(job => `
            <tr>
                <td>${this.escapeHtml(job.metadata.name)}</td>
                <td>${this.escapeHtml(job.metadata.namespace)}</td>
                <td>${this.escapeHtml(job.spec.schedule)}</td>
                <td>${job.spec.suspend ? 'Suspended' : 'Active'}</td>
            </tr>
        `).join('');
    }

    private showError(message: string): void {
        this.errorDiv.textContent = message;
        this.errorDiv.style.display = 'block';
    }

    private log(message: string): void {
        const time = new Date().toISOString();
        this.debugDiv.innerHTML += `<div>[${time}] ${message}</div>`;
        console.log(`[${time}] ${message}`);
    }

    private escapeHtml(unsafe: string): string {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
}

// Initialize the UI when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new CronJobUI();
}); 