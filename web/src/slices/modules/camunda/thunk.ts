import { createAsyncThunk } from '@reduxjs/toolkit';
import axiosInstance from '../../../utils/axiosInstance';

export const fetchResources = createAsyncThunk(
	'resources/fetchResources',
	async (
		{
			size,
			searchAfter,
			searchBefore,
			bpmnProcessId,
			startDate,
			endDate,
			processDefinitionKey,
			parentProcessInstanceKey,
			state,
		}: {
			size?: number;
			searchAfter?: string;
			searchBefore?: string;
			bpmnProcessId?: string;
			startDate?: string;
			endDate?: string;
			processDefinitionKey?: string;
			parentProcessInstanceKey?: string;
			state?: string;
		},
		{ rejectWithValue }
	) => {
		try {
			const params: Record<string, any> = {};
			if (size) params.size = size;
			if (searchAfter) params.searchAfter = searchAfter;
			if (searchBefore) params.searchBefore = searchBefore;
			if (bpmnProcessId) params.bpmnProcessId = bpmnProcessId;
			if (startDate) params.startDate = startDate;
			if (endDate) params.endDate = endDate;
			if (processDefinitionKey) params.processDefinitionKey = processDefinitionKey;
			if (parentProcessInstanceKey) params.parentProcessInstanceKey = parentProcessInstanceKey;
			if (state) params.state = state;

			const response = await axiosInstance.get('/camunda/process-instance', {
				params,
			});
			return response.data.data;
		} catch (error) {
			return rejectWithValue(error);
		}
	}
);
