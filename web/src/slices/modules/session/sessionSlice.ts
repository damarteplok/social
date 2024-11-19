import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { checkAuth } from './thunk';

interface Session {
	user: {
		id: string;
		name: string;
		email: string;
		image: string;
	};
}

interface SessionState {
	session: Session | null;
	loading?: boolean;
}

const initialState: SessionState = {
	session: null,
	loading: false,
};

const sessionSlice = createSlice({
	name: 'session',
	initialState,
	reducers: {
		setSession: (state, action: PayloadAction<Session | null>) => {
			state.session = action.payload;
		},
		resetSession: (state) => {
			state.session = null;
		},
	},
	extraReducers: (builder) => {
		builder
			.addCase(checkAuth.pending, (state) => {
				state.loading = true;
			})
			.addCase(checkAuth.fulfilled, (state, action) => {
				state.session = {
					user: {
						id: action.payload.id.toString(),
						name: action.payload.name,
						email: action.payload.email,
						image: state.session?.user.image || '',
					},
				};
				state.loading = false;
			})
			.addCase(checkAuth.rejected, (state) => {
				state.session = null;
				state.loading = false;
			});
	},
});

export const { setSession, resetSession } = sessionSlice.actions;
export default sessionSlice.reducer;
