import { combineReducers } from '@reduxjs/toolkit';
import sessionReducer from '../modules/session/sessionSlice';
import registerUserReducer from '../modules/auth/registerSlice';
import dashboardReducer from '../modules/dashboard/dashboardSlice';
import camundaReducer from '../modules/camunda/camundaSlice';
import forgotPasswordReducer from '../modules/forgotPassword/forgotPasswordSlice';

const appReducer = combineReducers({
	session: sessionReducer,
	registerUser: registerUserReducer,
	dashboard: dashboardReducer,
	camunda: camundaReducer,
	forgotPassword: forgotPasswordReducer,
});

const rootReducer = (state: any, action: any) => {
	if (action.type === 'RESET_STORE') {
		state = undefined;
	}
	return appReducer(state, action);
};

export type RootState = ReturnType<typeof rootReducer>;

export default rootReducer;
