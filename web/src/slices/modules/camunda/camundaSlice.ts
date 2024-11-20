import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { fetchResources } from './thunk';

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
	loading: boolean;
	error: string | null;
}

const initialState: CamundaState = {
	resources: null,
	loading: false,
	error: null,
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
			});
	},
});

export const { setResources, resetResources } = camundaSlice.actions;
export default camundaSlice.reducer;
