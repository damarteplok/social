import { createAsyncThunk } from '@reduxjs/toolkit';
import axiosInstance from '../../../utils/axiosInstance';

export const fetchDashboard = createAsyncThunk(
	'dashboard/fetchDashboard',
	async (_, { rejectWithValue }) => {
		try {
			const response = await axiosInstance.get('/health');
			const responseOperate = await axiosInstance.get('/camunda/resource/operate/statistics');
			const stats = responseOperate.data.data.stats;
			const process = responseOperate.data.data.process;
			return {
				status: response.data.data.status,
				env: response.data.data.env,
				version: response.data.data.version,
				active: stats.active,
				running: stats.running,
				incident: stats.withIncidents,
				processStats: process,
			}
		} catch (error) {
			return rejectWithValue(error);
		}
	}
);
