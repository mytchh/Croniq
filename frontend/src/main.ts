import { CronJob, CronJobResponse, ClusterInfo, JobStats, JobResponse, Job } from './types';

class CronJobUI {
    private tableBody: HTMLElement;
    private errorDiv: HTMLElement;
    private clusterInfoDiv: HTMLElement;
    private debugDiv: HTMLElement;
    private activeTab: string = 'dashboard';

    constructor() {
        this.tableBody = document.getElementById('cronJobList') as HTMLElement;
        this.errorDiv = document.getElementById('error') as HTMLElement;
        this.clusterInfoDiv = document.getElementById('clusterInfo') as HTMLElement;
        this.debugDiv = document.getElementById('debug') as HTMLElement;
        this.initTabs();
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
            const response = await fetch('http://localhost:8080/api/cluster-info', {
                headers: {
                    'Cache-Control': 'no-cache',
                    'Pragma': 'no-cache'
                }
            });
            console.log('Cluster info response:', response.status);
            if (!response.ok) {
                const text = await response.text();
                console.error('Error response:', text);
                throw new Error('Failed to fetch cluster info');
            }
            const info: ClusterInfo = await response.json();
            console.log('Cluster info received:', info);
            this.renderClusterInfo(info);
        } catch (err) {
            console.error('Fetch error:', err);
            this.showError(err instanceof Error ? err.message : 'An error occurred');
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
            const response = await fetch('http://localhost:8080/api/cronjobs', {
                headers: {
                    'Cache-Control': 'no-cache',
                    'Pragma': 'no-cache'
                }
            });
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

    private initTabs(): void {
        const tabs = document.querySelectorAll('.tab-button');
        tabs.forEach(tab => {
            tab.addEventListener('click', (e) => {
                const target = e.target as HTMLElement;
                const tabName = target.dataset.tab;
                if (tabName) {
                    this.switchTab(tabName);
                }
            });
        });
    }

    private switchTab(tabName: string): void {
        // Update buttons
        document.querySelectorAll('.tab-button').forEach(tab => {
            tab.classList.toggle('active', tab.getAttribute('data-tab') === tabName);
        });

        // Update content
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.toggle('active', content.id === tabName);
        });

        this.activeTab = tabName;
        this.refreshActiveTab();
    }

    private async refreshActiveTab(): Promise<void> {
        switch (this.activeTab) {
            case 'dashboard':
                await this.fetchClusterInfo();
                await this.fetchStats();
                break;
            case 'cronjobs':
                await this.fetchAndDisplayCronJobs();
                break;
            case 'jobs':
                await this.fetchAndDisplayJobs();
                break;
        }
    }

    private async fetchStats(): Promise<void> {
        try {
            const response = await fetch('http://localhost:8080/api/stats', {
                headers: {
                    'Cache-Control': 'no-cache',
                    'Pragma': 'no-cache'
                }
            });
            
            if (!response.ok) {
                throw new Error('Failed to fetch stats');
            }
            
            const stats: JobStats = await response.json();
            this.renderStats(stats);
        } catch (err) {
            this.showError(err instanceof Error ? err.message : 'An error occurred');
        }
    }

    private renderStats(stats: JobStats): void {
        const cronJobStatsDiv = document.getElementById('cronJobStats');
        const jobStatsDiv = document.getElementById('jobStats');
        
        if (cronJobStatsDiv) {
            cronJobStatsDiv.innerHTML = `
                <h3>CronJobs</h3>
                <p>Total: ${stats.totalCronJobs}</p>
                <p>Active: ${stats.activeCronJobs}</p>
            `;
        }
        
        if (jobStatsDiv) {
            jobStatsDiv.innerHTML = `
                <h3>Jobs</h3>
                <p>Total: ${stats.totalJobs}</p>
                <p>Running: ${stats.runningJobs}</p>
                <p>Succeeded: ${stats.succeededJobs}</p>
                <p>Failed: ${stats.failedJobs}</p>
            `;
        }
    }

    private async fetchAndDisplayJobs(): Promise<void> {
        try {
            const response = await fetch('http://localhost:8080/api/jobs', {
                headers: {
                    'Cache-Control': 'no-cache',
                    'Pragma': 'no-cache'
                }
            });
            
            if (!response.ok) {
                throw new Error('Failed to fetch jobs');
            }
            
            const data: JobResponse = await response.json();
            this.renderJobs(data.items || []);
        } catch (err) {
            this.showError(err instanceof Error ? err.message : 'An error occurred');
        }
    }

    private renderJobs(jobs: Job[]): void {
        const jobList = document.getElementById('jobList');
        if (!jobList) return;

        jobList.innerHTML = jobs.map(job => `
            <tr>
                <td>${this.escapeHtml(job.metadata.name)}</td>
                <td>${this.escapeHtml(job.metadata.namespace)}</td>
                <td>${this.getJobStatus(job)}</td>
                <td>${job.status.startTime ? new Date(job.status.startTime).toLocaleString() : 'N/A'}</td>
                <td>${job.status.completionTime ? new Date(job.status.completionTime).toLocaleString() : 'N/A'}</td>
            </tr>
        `).join('');
    }

    private getJobStatus(job: Job): string {
        if (job.status.active > 0) return 'Running';
        if (job.status.succeeded > 0) return 'Succeeded';
        if (job.status.failed > 0) return 'Failed';
        return 'Pending';
    }
}

// Initialize the UI when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new CronJobUI();
}); 