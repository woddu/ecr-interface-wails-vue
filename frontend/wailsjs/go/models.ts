export namespace main {
	
	export class ScoresResult {
	    wwHighestScores: number[];
	    ptHighestScores: number[];
	    examHighestScore: number;
	    weightedScores: number[];
	
	    static createFrom(source: any = {}) {
	        return new ScoresResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.wwHighestScores = source["wwHighestScores"];
	        this.ptHighestScores = source["ptHighestScores"];
	        this.examHighestScore = source["examHighestScore"];
	        this.weightedScores = source["weightedScores"];
	    }
	}

}

