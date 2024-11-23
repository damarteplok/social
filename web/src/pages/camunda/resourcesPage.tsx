import React, { useEffect, useState } from 'react';
import { Grid2, Fab } from '@mui/material';
import TableComponent from '../../components/TableComponent';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../slices/store/rootReducer';
import { fetchResources } from '../../slices/modules/camunda/thunk';
import { PageContainer, useNotifications } from '@toolpad/core';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import ScrollToTop from 'react-scroll-to-top';
import { columnsResources } from './filters/columnsResource';
import FilterResourceComponent from './filters/FilterResourceComponent';
import PageToolbarResource from './filters/PageToolbarResource';

const ResourcesPage: React.FC = () => {
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
	const [rowsPerPage, setRowsPerPage] = useState(10);
	const [searchAfter, setSearchAfter] = useState('');
	const [searchBefore, setSearchBefore] = useState('');
	const [filter, setFilter] = useState(initFilter);

	const { loading, resources, error } = useSelector(
		(state: RootState) => state.camunda
	);

	const onNextPage = () => {
		setSearchAfter(
			resources?.items[resources.items.length - 1]?.key?.toString() || ''
		);
		setSearchBefore('');
	};

	const onPrevPage = () => {
		setSearchAfter('');
		setSearchBefore(resources?.items[0]?.key?.toString() || '');
	};

	const handleView = (id: string) => {
		console.log('View', id);
	};

	const handleEdit = (id: string) => {
		console.log('Edit', id);
	};

	const handleDelete = (id: string) => {};

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

	//fetch first fetchResources
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
					totalRows={resources?.total || 0}
					onPageChange={setPage}
					onRowsPerPageChange={setRowsPerPage}
					showActions={false}
					showView={true}
					showEdit={true}
					showDelete={true}
					showNumber={true}
					filter={filter}
					initFilter={initFilter}
					onNextPage={onNextPage}
					onPrevPage={onPrevPage}
					loading={loading}
				/>
			</Grid2>
			<ScrollToTop
				smooth
				component={
					<Fab size='small' aria-label='scroll back to top'>
						<KeyboardArrowUpIcon />
					</Fab>
				}
			/>
		</PageContainer>
	);
};

export default ResourcesPage;
