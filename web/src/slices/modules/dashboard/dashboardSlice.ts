import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { fetchDashboard } from './thunk';

interface ProcessType {
	processId: string;
	version: number;
	name: string;
	bpmnProcessId: string;
	instancesWithActiveIncidentsCount: number;
	activeInstancesCount: number;
}

interface ProcessStats {
	bpmnProcessId: string;
	processName: string;
	instancesWithActiveIncidentsCount: number;
	activeInstancesCount: number;
	processes: ProcessType[];
}

interface Dashboard {
	dashboard: {
		health: {
			status: string;
			env: string;
			version: string;
		};
		count: {
			active: number;
			running: number;
			incident: number;
		};
		processStats: ProcessStats[];
	};
}

interface DashboardState {
	dashboard: Dashboard | null;
	loading: boolean;
	error: string | null;
}

const initialState: DashboardState = {
	dashboard: null,
	loading: false,
	error: null,
};

const dashboardSlice = createSlice({
	name: 'dashboard',
	initialState,
	reducers: {
		setDashboard: (state, action: PayloadAction<Dashboard | null>) => {
			state.dashboard = action.payload;
		},
		resetDashboard: (state) => {
			state.dashboard = null;
		},
	},
	extraReducers: (builder) => {
		builder
			.addCase(fetchDashboard.pending, (state) => {
				state.loading = true;
			})
			.addCase(fetchDashboard.fulfilled, (state, action) => {
				state.dashboard = {
					dashboard: {
						health: {
							status: action.payload.status,
							env: action.payload.env,
							version: action.payload.version,
						},
						count: {
							active: action.payload.active,
							running: action.payload.running,
							incident: action.payload.incident,
						},
						processStats: action.payload.processStats,
					},
				};
				state.loading = false;
				state.error = null;
			})
			.addCase(fetchDashboard.rejected, (state, action) => {
				state.loading = false;
				state.error = action.payload as string;
			});
	},
});

export const { setDashboard, resetDashboard } = dashboardSlice.actions;
export default dashboardSlice.reducer;
