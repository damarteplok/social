import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { forgotPassword } from './thunk';

interface ForgotPassword {
	email: string;
}

interface ForgotPasswordState {
	loading: boolean;
	error?: string;
	success?: string;
	forgotPassword: ForgotPassword | null;
}

const initialState: ForgotPasswordState = {
	loading: false,
	error: '',
	success: '',
	forgotPassword: null,
};

const forgotPasswordSlice = createSlice({
	name: 'forgotPassword',
	initialState,
	reducers: {
		setForgotPassword: (
			state,
			action: PayloadAction<ForgotPassword | null>
		) => {
			state.forgotPassword = action.payload;
		},
		resetForgotPassword: (state) => {
			state.forgotPassword = null;
		},
	},
	extraReducers(builder) {
		builder
			.addCase(forgotPassword.pending, (state) => {
				state.loading = true;
				state.error = '';
				state.success = '';
			})
			.addCase(forgotPassword.fulfilled, (state, action) => {
				state.forgotPassword = action.payload;
				state.loading = false;
				state.error = '';
				state.success = 'Success';
			})
			.addCase(forgotPassword.rejected, (state, action) => {
				state.loading = false;
				state.error = (action.payload as string) || 'Error';
				state.success = '';
			});
	},
});

export const { setForgotPassword, resetForgotPassword } =
	forgotPasswordSlice.actions;
export default forgotPasswordSlice.reducer;
