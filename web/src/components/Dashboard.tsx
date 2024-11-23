import {
	Grid2,
	Card,
	CardContent,
	Typography,
	Skeleton,
	Accordion,
	AccordionDetails,
	AccordionSummary,
} from '@mui/material';
import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { fetchDashboard } from '../slices/modules/dashboard/thunk';
import { RootState } from '../slices/store/rootReducer';
import { useNotifications } from '@toolpad/core/useNotifications';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';

const Dashboard = () => {
	const dispatch = useDispatch<any>();
	const notifications = useNotifications();

	const { dashboard, loading, error } = useSelector(
		(state: RootState) => state.dashboard
	);

	useEffect(() => {
		dispatch(fetchDashboard());
	}, [dispatch]);

	useEffect(() => {
		if (error) {
			notifications.show(error, {
				severity: 'error',
				autoHideDuration: 3000,
			});
		}
	}, [error, notifications]);

	return (
		<Grid2 container spacing={3}>
			<Grid2 size={{ xs: 12, sm: 12, md: 12 }}>
				<Card
					variant='outlined'
					sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}
				>
					<CardContent sx={{ flexGrow: 1 }}>
						{loading ? (
							<>
								<Skeleton variant='text' />
								<Skeleton variant='rectangular' height={160} />
							</>
						) : (
							<>
								<Typography variant='h5'>Health</Typography>
								<Typography variant='body2'>
									env: {dashboard?.dashboard.health.env}
									<br />
									status : {dashboard?.dashboard.health.status}
									<br />
									version : {dashboard?.dashboard.health.version}
								</Typography>
							</>
						)}
					</CardContent>
				</Card>
			</Grid2>
			<Grid2 size={{ xs: 12, sm: 6, md: 4 }}>
				<Card
					variant='outlined'
					sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}
				>
					<CardContent sx={{ flexGrow: 1 }}>
						{loading ? (
							<>
								<Skeleton variant='text' />
								<Skeleton variant='rectangular' height={100} />
							</>
						) : (
							<>
								<Typography variant='h5'>Running</Typography>
								<Typography variant='h4'>
									{dashboard?.dashboard.count.running}
								</Typography>
							</>
						)}
					</CardContent>
				</Card>
			</Grid2>
			<Grid2 size={{ xs: 12, sm: 6, md: 4 }}>
				<Card
					variant='outlined'
					sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}
				>
					<CardContent sx={{ flexGrow: 1 }}>
						{loading ? (
							<>
								<Skeleton variant='text' />
								<Skeleton variant='rectangular' height={100} />
							</>
						) : (
							<>
								<Typography variant='h5'>Active</Typography>
								<Typography variant='h4'>
									{dashboard?.dashboard.count.active}
								</Typography>
							</>
						)}
					</CardContent>
				</Card>
			</Grid2>
			<Grid2 size={{ xs: 12, sm: 6, md: 4 }}>
				<Card
					variant='outlined'
					sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}
				>
					<CardContent sx={{ flexGrow: 1 }}>
						{loading ? (
							<>
								<Skeleton variant='text' />
								<Skeleton variant='rectangular' height={100} />
							</>
						) : (
							<>
								<Typography variant='h5'>Incidents</Typography>
								<Typography variant='h4'>
									{dashboard?.dashboard.count.incident}
								</Typography>
							</>
						)}
					</CardContent>
				</Card>
			</Grid2>

			<Grid2 size={{ xs: 12, sm: 12, md: 12 }}>
				<Card
					variant='outlined'
					sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}
				>
					<CardContent sx={{ flexGrow: 1 }}>
						{loading ? (
							<>
								<Skeleton variant='text' />
								<Skeleton variant='rectangular' height={200} />
							</>
						) : (
							<>
								<Typography variant='h5' sx={{mb: 2}}>Process Instance by Name</Typography>

								{dashboard?.dashboard.processStats.map((process) => (
									<Accordion key={process.bpmnProcessId}>
										<AccordionSummary
											expandIcon={<ExpandMoreIcon />}
											aria-controls='panel1-content'
											id={process.bpmnProcessId}
										>
											{process.processName}
										</AccordionSummary>
										<AccordionDetails>
											<Typography variant='body2' sx={{mb: 2}}>
												active: {process.activeInstancesCount} incidents:{' '}
												{process.instancesWithActiveIncidentsCount}
											</Typography>
											{process.processes.map((subProcess) => (
												<Accordion key={subProcess.processId}>
													<AccordionSummary
														expandIcon={<ExpandMoreIcon />}
														aria-controls='panel2-content'
														id={subProcess.processId}
													>
														{subProcess.name} (Version: {subProcess.version})
													</AccordionSummary>
													<AccordionDetails>
														<Typography variant='body2'>
															active: {subProcess.activeInstancesCount}{' '}
															<br />
															incidents:{' '}
															{subProcess.instancesWithActiveIncidentsCount}
														</Typography>
													</AccordionDetails>
												</Accordion>
											))}
										</AccordionDetails>
									</Accordion>
								))}
							</>
						)}
					</CardContent>
				</Card>
			</Grid2>
		</Grid2>
	);
};
export default Dashboard;
