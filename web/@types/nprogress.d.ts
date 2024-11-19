declare module 'nprogress' {
	interface NProgress {
		start: () => NProgress;
		done: (force?: boolean) => NProgress;
		set: (n: number) => NProgress;
		inc: (amount?: number) => NProgress;
		configure: (options: Partial<NProgressOptions>) => NProgress;
		status: null | number;
		isStarted: () => boolean;
	}

	interface NProgressOptions {
		minimum: number;
		easing: string;
		positionUsing: string;
		speed: number;
		trickle: boolean;
		trickleRate: number;
		trickleSpeed: number;
		showSpinner: boolean;
		barSelector: string;
		spinnerSelector: string;
		parent: string;
		template: string;
	}

	const nprogress: NProgress;
	export default nprogress;
}
