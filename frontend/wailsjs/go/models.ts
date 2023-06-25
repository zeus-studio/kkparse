export namespace parser {
	
	export class VideoParseInfo {
	    // Go type: struct { Uid string "json:\"uid\""; Name string "json:\"name\""; Avatar string "json:\"avatar\"" }
	    author: any;
	    title: string;
	    video_url: string;
	    music_url: string;
	    cover_url: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new VideoParseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.author = this.convertValues(source["author"], Object);
	        this.title = source["title"];
	        this.video_url = source["video_url"];
	        this.music_url = source["music_url"];
	        this.cover_url = source["cover_url"];
	        this.source = source["source"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

