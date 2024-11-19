import * as React from 'react';
import DashboardIcon from '@mui/icons-material/Dashboard';
import { Outlet, useNavigate } from 'react-router-dom';
import { setAuthToken } from './utils/axiosInstance';
import { useDispatch, useSelector } from 'react-redux';
import { setSession } from './slices/modules/session/sessionSlice';
import { createTheme } from '@mui/material/styles';
import { AppProvider } from '@toolpad/core/react-router-dom';
import CloudIcon from '@mui/icons-material/Cloud';
import AssignmentIndIcon from '@mui/icons-material/AssignmentInd';
import AssignmentTurnedInIcon from '@mui/icons-material/AssignmentTurnedIn';
import MonitorHeartIcon from '@mui/icons-material/MonitorHeart';
import FeedbackIcon from '@mui/icons-material/Feedback';
import { SessionContext, type Navigation } from '@toolpad/core/AppProvider';
import { RootState } from './slices/store/rootReducer';

const NAVIGATION: Navigation = [
	{
		title: 'Dashboard',
		icon: <MonitorHeartIcon />,
	},
	{
		kind: 'divider',
	},
	{
		kind: 'header',
		title: 'Camunda Rest API',
	},
	{
		segment: 'camunda',
		title: 'Camunda',
		icon: <CloudIcon />,
		children: [
			{
				segment: 'resources',
				title: 'Resources',
				icon: <DashboardIcon />,
				pattern: 'resources',
			},
			{
				segment: 'process-instance',
				title: 'Process Instance',
				icon: <AssignmentTurnedInIcon />,
				pattern: 'process-instance',
			},
			{
				segment: 'user-task',
				title: 'User Task',
				icon: <AssignmentIndIcon />,
				pattern: 'user-task',
			},
			{
				segment: 'incidents',
				title: 'Incidents',
				icon: <FeedbackIcon />,
				pattern: 'incident',
			},
		],
	},
];

const customTheme = createTheme({
	cssVariables: {
		colorSchemeSelector: 'data-toolpad-color-scheme',
	},
	colorSchemes: { light: true, dark: true },
	breakpoints: {
		values: {
			xs: 0,
			sm: 600,
			md: 600,
			lg: 1200,
			xl: 1536,
		},
	},
});

const BRANDING = {
	title: (import.meta.env.VITE_APP_TITLE as string) || 'Damarmunda',
};

export default function App() {
	const session = useSelector((state: RootState) => state.session.session);
	const dispatch = useDispatch();

	const navigate = useNavigate();

	const signIn = React.useCallback(() => {
		navigate('/sign-in');
	}, [navigate]);

	const signOut = React.useCallback(() => {
		dispatch(setSession(null));
		setAuthToken(null);
		navigate('/sign-in');
	}, [dispatch, navigate]);

	return (
		<SessionContext.Provider value={session}>
			<AppProvider
				navigation={NAVIGATION}
				branding={BRANDING}
				session={session}
				authentication={{ signIn, signOut }}
				theme={customTheme}
			>
				<Outlet />
			</AppProvider>
		</SessionContext.Provider>
	);
}
