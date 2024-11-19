import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import Layout from './layouts/dashboard';
import OrdersPage from './pages/orders.tsx';
import SignInPage from './pages/signIn';
import { Provider } from 'react-redux';
import { PersistGate } from 'redux-persist/integration/react';
import { store, persistor } from './slices/store';
import 'nprogress/nprogress.css';
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import PrivateRoute from './components/PrivateRoute';
import ForgotPasswordPage from './pages/forgotPasswordPage.tsx';
import SignUpPage from './pages/signUpPage.tsx';
import DashboardPage from './pages/dashboardPage.tsx';
import NotFoundPage from './pages/notFoundPage.tsx';
import UnauthorizedPage from './pages/unAuthorizedPage.tsx';
import ResourcesPage from './pages/camunda/resourcesPage.tsx';

const router = createBrowserRouter([
	{
		Component: App,
		children: [
			{
				path: '/',
				Component: Layout,
				children: [
					{
						path: '',
						Component: () => (
							<PrivateRoute>
								<DashboardPage />
							</PrivateRoute>
						),
					},
					{
						path: 'camunda',
						children: [
							{
								path: 'resources',
								Component: () => (
									<PrivateRoute>
										<ResourcesPage />
									</PrivateRoute>
								),
							},
							{
								path: 'process-instance',
								Component: () => (
									<PrivateRoute>
										<OrdersPage />
									</PrivateRoute>
								),
							},
							{
								path: 'user-task',
								Component: () => (
									<PrivateRoute>
										<OrdersPage />
									</PrivateRoute>
								),
							},
						],
					},
				],
			},
			{
				path: '/sign-in',
				Component: SignInPage,
			},
			{
				path: '/sign-up',
				Component: SignUpPage,
			},
			{
				path: '/forgot-password',
				Component: ForgotPasswordPage,
			},
			{
				path: '/unauthorized', // Add route for UnauthorizedPage
				element: <UnauthorizedPage />,
			},
			{
				path: '*',
				Component: NotFoundPage,
			},
		],
	},
]);

createRoot(document.getElementById('root')!).render(
	<StrictMode>
		<Provider store={store}>
			<PersistGate loading={null} persistor={persistor}>
				<RouterProvider router={router} />
			</PersistGate>
		</Provider>
	</StrictMode>
);
