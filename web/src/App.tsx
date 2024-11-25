import * as React from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import { setAuthToken } from './utils/axiosInstance';
import { useDispatch, useSelector } from 'react-redux';
import { setSession } from './slices/modules/session/sessionSlice';
import { createTheme } from '@mui/material/styles';
import { AppProvider } from '@toolpad/core/react-router-dom';
import CloudIcon from '@mui/icons-material/Cloud';
import MonitorHeartIcon from '@mui/icons-material/MonitorHeart';
import { SessionContext, type Navigation } from '@toolpad/core/AppProvider';
import { RootState } from './slices/store/rootReducer';
import HomeIcon from '@mui/icons-material/Home';

const NAVIGATION: Navigation = [
	{
		title: 'Home',
		icon: <HomeIcon />,
	},
	{
		kind: 'divider',
	},
	{
		kind: 'header',
		title: 'Monitoring',
	},
	{
		segment: 'monitoring',
		title: 'Monitoring',
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
	},
	// {
	// 	segment: 'camunda',
	// 	title: 'Camunda',
	// 	icon: <CloudIcon />,
	// 	children: [
	// 		{
	// 			segment: 'resources',
	// 			title: 'Resources',
	// 			icon: <DashboardIcon />,
	// 			pattern: 'resources',
	// 		},
	// 	],
	// },
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
