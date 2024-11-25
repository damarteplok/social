import React, { useEffect, useState } from 'react';
import {
	Box,
	Button,
	Grid2,
	IconButton,
	Modal,
	Skeleton,
	TextField,
	Typography,
} from '@mui/material';
import TableComponent from '../../components/view/table/TableComponent';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../slices/store/rootReducer';
import {
	fetchBpmnXml,
	fetchResources,
	resolveIncident,
} from '../../slices/modules/camunda/thunk';
import { PageContainer, useNotifications } from '@toolpad/core';
import { columnsResources, RowResource } from './filters/columnsResource';
import FilterResourceComponent from './filters/FilterResourceComponent';
import PageToolbarResource from './filters/PageToolbarResource';
import ScrollToTopButton from '../../components/view/scroll/ScrollToTopButton';
import BpmnViewerComponent from '../../components/view/bpmn/BpmnViewerComponent';
import { modalStyle } from '../../layouts/modal';
import ReplayIcon from '@mui/icons-material/Replay';

const CamundaPage: React.FC = () => {
	const notifications = useNotifications();
	const dispatch = useDispatch<any>();

	const initFilter = {
		bpmnProcessId: '',
		startDate: '',
		endDate: '',
		processDefinitionKey: '',
		parentProcessInstanceKey: '',
		state: '',
	};

	const [page, setPage] = useState(0);
	const [rowDetail, setRowDetail] = useState<RowResource | null>(null);
	const [rowsPerPage, setRowsPerPage] = useState(10);
	const [searchAfter, setSearchAfter] = useState('');
	const [searchBefore, setSearchBefore] = useState('');
	const [filter, setFilter] = useState(initFilter);
	const [open, setOpen] = useState(false);

	const {
		loading,
		resources,
		error,
		loadingViewer,
		errorViewer,
		bpmnXml,
		errorIncident,
		loadingIncident,
		successIncident,
	} = useSelector((state: RootState) => state.camunda);

	const onNextPage = () => {
		setSearchAfter(
			resources?.items[resources.items.length - 1]?.key?.toString() ?? ''
		);
		setSearchBefore('');
	};

	const onPrevPage = () => {
		setSearchAfter('');
		setSearchBefore(resources?.items[0]?.key?.toString() ?? '');
	};

	const handleView = (row: RowResource) => {
		const processDefinitionKey = row.processDefinitionKey;
		dispatch(fetchBpmnXml(processDefinitionKey.toString()));
		setRowDetail(row);
		setOpen(true);
	};

	const handleCustomAction = (row: RowResource) => {
		dispatch(resolveIncident(row.key));
	};

	const handleEdit = (row: RowResource) => {};

	const handleDelete = (row: RowResource) => {};

	const handleApplyFilter = (filter: any) => {
		setPage(0);
		setFilter(filter);
		dispatch(
			fetchResources({
				size: rowsPerPage,
				searchAfter,
				searchBefore,
				...filter,
			})
		);
	};

	const fetchData = () => {
		dispatch(
			fetchResources({
				size: rowsPerPage,
				searchAfter,
				searchBefore,
				...filter,
			})
		);
	};

	useEffect(() => {
		fetchData();
	}, [dispatch, rowsPerPage, searchAfter, searchBefore]);

	useEffect(() => {
		if (error) {
			notifications.show(error, {
				severity: 'error',
				autoHideDuration: 3000,
			});
		}
	}, [error, notifications]);

	useEffect(() => {
		if (errorViewer) {
			notifications.show(errorViewer, {
				severity: 'error',
				autoHideDuration: 3000,
			});
		}
	}, [errorViewer, notifications]);

	useEffect(() => {
		if (errorIncident) {
			notifications.show(errorIncident, {
				severity: 'error',
				autoHideDuration: 3000,
			});
		}
		if (successIncident) {
			notifications.show('Incident resolved', {
				severity: 'success',
				autoHideDuration: 3000,
			});
		}
	}, [errorIncident, successIncident, notifications]);

	return (
		<PageContainer slots={{ toolbar: PageToolbarResource }}>
			<Grid2 container spacing={3}>
				<TableComponent
					columns={columnsResources}
					rows={resources?.items || []}
					onView={handleView}
					onEdit={handleEdit}
					onDelete={handleDelete}
					onApplyFilter={handleApplyFilter}
					showFilter
					showSearch={false}
					customContent={
						<FilterResourceComponent filter={filter} setFilter={setFilter} />
					}
					page={page}
					rowsPerPage={rowsPerPage}
					totalRows={resources?.total ?? 0}
					onPageChange={setPage}
					onRowsPerPageChange={setRowsPerPage}
					showActions={true}
					showView={true}
					showEdit={false}
					showDelete={true}
					showNumber={true}
					filter={filter}
					initFilter={initFilter}
					onNextPage={onNextPage}
					onPrevPage={onPrevPage}
					loading={loading}
					customActions={(row: RowResource) => {
						return row.incident ? (
							<IconButton
								onClick={() => handleCustomAction(row)}
								disabled={loadingIncident}
							>
								<ReplayIcon />
							</IconButton>
						) : (
							<></>
						);
					}}
				/>
			</Grid2>
			<ScrollToTopButton />
			<Modal
				open={open}
				onClose={() => {
					setOpen(false);
					setRowDetail(null);
				}}
			>
				<Box sx={{ ...modalStyle }}>
					{loadingViewer && (
						<>
							<Skeleton
								variant='rectangular'
								width='100%'
								height={280}
								sx={{ mb: 1 }}
							/>
							{Array.from({ length: 3 }).map((_, index) => (
								<Skeleton
									key={'skeleton' + index}
									variant='text'
									width='100%'
									height={40}
									sx={{ mb: 1 }}
								/>
							))}
						</>
					)}
					{!loadingViewer && errorViewer && <Box>{errorViewer}</Box>}
					{!loadingViewer && !errorViewer && !bpmnXml && <Box>No Data</Box>}
					{!loadingViewer && !errorViewer && bpmnXml && (
						<>
							<Typography variant='h5' sx={{ mb: 2 }}>
								{rowDetail?.bpmnProcessId} {rowDetail?.processDefinitionKey}
							</Typography>
							<BpmnViewerComponent xml={bpmnXml} />
						</>
					)}
				</Box>
			</Modal>
		</PageContainer>
	);
};

export default CamundaPage;
