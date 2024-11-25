import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { fetchBpmnXml, fetchResources, resolveIncident } from './thunk';

interface ResourceCamundaItem {
	bpmnProcessId: string;
	endDate: string;
	incident: boolean;
	key: number;
	processDefinitionKey: number;
	processVersion: number;
	startDate: string;
	state: string;
}

interface ResourceCamunda {
	items: Array<ResourceCamundaItem>;
	sortValues: Array<string>;
	total: number;
}

interface CamundaState {
	resources: ResourceCamunda | null;
	bpmnXml: string | null;
	loading: boolean;
	error: string | null;
	loadingViewer: boolean;
	errorViewer: string | null;
	loadingIncident: boolean;
	errorIncident: string | null;
	successIncident: boolean;
}

const initialState: CamundaState = {
	resources: null,
	loading: false,
	error: null,
	loadingViewer: false,
	errorViewer: null,
	bpmnXml: null,
	loadingIncident: false,
	errorIncident: null,
	successIncident: false,
};

export const camundaSlice = createSlice({
	name: 'camunda',
	initialState,
	reducers: {
		setResources: (state, action: PayloadAction<ResourceCamunda | null>) => {
			state.resources = action.payload;
		},
		resetResources: (state) => {
			state.resources = null;
		},
	},
	extraReducers: (builder) => {
		builder
			.addCase(fetchResources.pending, (state) => {
				state.loading = true;
				state.error = null;
			})
			.addCase(fetchResources.fulfilled, (state, action) => {
				state.resources = action.payload;
				state.loading = false;
				state.error = null;
			})
			.addCase(fetchResources.rejected, (state, action) => {
				state.resources = null;
				state.loading = false;
				state.error = action.payload as string;
			})
			.addCase(fetchBpmnXml.pending, (state) => {
				state.loadingViewer = true;
				state.errorViewer = null;
			})
			.addCase(fetchBpmnXml.fulfilled, (state, action) => {
				state.bpmnXml = action.payload;
				state.loadingViewer = false;
				state.errorViewer = null;
			})
			.addCase(fetchBpmnXml.rejected, (state, action) => {
				state.bpmnXml = null;
				state.loadingViewer = false;
				state.errorViewer = action.payload as string;
			})
			.addCase(resolveIncident.pending, (state) => {
				state.loadingIncident = true;
				state.errorIncident = null;
				state.successIncident = false;
			})
			.addCase(resolveIncident.fulfilled, (state) => {
				state.loadingIncident = false;
				state.errorIncident = null;
				state.successIncident = true;
			})
			.addCase(resolveIncident.rejected, (state, action) => {
				state.loadingIncident = false;
				state.errorIncident = action.payload as string;
				state.successIncident = false;
			});
	},
});

export const { setResources, resetResources } = camundaSlice.actions;
export default camundaSlice.reducer;
