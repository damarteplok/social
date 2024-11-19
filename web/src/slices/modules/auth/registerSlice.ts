import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { registerUser } from './thunk';

interface RegisterUser {
	username: string;
	email: string;
	password: string;
}

interface RegisterUserState {
	registerUser: RegisterUser | null;
	loading?: boolean;
	success?: boolean;
	errorMessage?: string;
}

const initialState: RegisterUserState = {
	registerUser: null,
	loading: false,
	success: false,
	errorMessage: '',
};

const registerUserSlice = createSlice({
	name: 'registerUser',
	initialState,
	reducers: {
		setRegisterUser: (state, action: PayloadAction<RegisterUser | null>) => {
			state.registerUser = action.payload;
		},
		resetRegisterUser: (state) => {
			state.registerUser = null;
		},
	},
	extraReducers(builder) {
		builder
			.addCase(registerUser.pending, (state) => {
				state.loading = true;
			})
			.addCase(registerUser.fulfilled, (state, action) => {
				state.registerUser = action.payload;
				state.loading = false;
				state.success = true;
				state.errorMessage = '';
			})
			.addCase(registerUser.rejected, (state, action) => {
				state.loading = false;
				state.success = false;
				state.errorMessage = action.payload as string;
			});
	},
});

export const { setRegisterUser, resetRegisterUser } = registerUserSlice.actions;
export default registerUserSlice.reducer;
