export namespace models {
	
	export class HostGroup {
	    id: string;
	    title: string;
	    content: string;
	    isActive: boolean;
	    isRemote: boolean;
	    remoteUrl?: string;
	    lastUpdated: string;
	    createdAt: string;
	    tags?: string[];
	    parent?: string;
	
	    static createFrom(source: any = {}) {
	        return new HostGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.content = source["content"];
	        this.isActive = source["isActive"];
	        this.isRemote = source["isRemote"];
	        this.remoteUrl = source["remoteUrl"];
	        this.lastUpdated = source["lastUpdated"];
	        this.createdAt = source["createdAt"];
	        this.tags = source["tags"];
	        this.parent = source["parent"];
	    }
	}
	export class HostsFile {
	    path: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new HostsFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.content = source["content"];
	    }
	}
	export class SystemConfig {
	    theme: string;
	    autoBackup: boolean;
	    backupFolder: string;
	    updateFreq: number;
	    hostsMode: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.autoBackup = source["autoBackup"];
	        this.backupFolder = source["backupFolder"];
	        this.updateFreq = source["updateFreq"];
	        this.hostsMode = source["hostsMode"];
	    }
	}

}

