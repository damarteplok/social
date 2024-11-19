import React, { useEffect, useState } from 'react';
import { Box, Grid2, Skeleton, TextField } from '@mui/material';
import TableComponent from '../../components/TableComponent';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../slices/store/rootReducer';
import { fetchResources } from '../../slices/modules/camunda/thunk';

const columns = [
	{ id: 'bpmnProcessId', label: 'Bpmn Process Id', minWidth: 170 },
	{ id: 'key', label: 'Key', minWidth: 170 },
	{
		id: 'processDefinitionKey',
		label: 'Process Definition Key',
		minWidth: 170,
	},
	{ id: 'processVersion', label: 'Version', minWidth: 170 },
	{ id: 'startDate', label: 'Start Date', minWidth: 170 },
	{ id: 'endDate', label: 'End Date', minWidth: 170 },
	{ id: 'incident', label: 'Incident', minWidth: 170 },
	{ id: 'state', label: 'State', minWidth: 170 },
];

const ResourcesPage: React.FC = () => {
	const dispatch = useDispatch<any>();
	const [page, setPage] = useState(0);
	const [rowsPerPage, setRowsPerPage] = useState(10);
	const [searchAfter, setSearchAfter] = useState('');
	const [searchBefore, setSearchBefore] = useState('');

	const { loading, resources } = useSelector(
		(state: RootState) => state.camunda
	);

	const handleView = (id: string) => {
		console.log('View', id);
	};

	const handleEdit = (id: string) => {
		console.log('Edit', id);
	};

	const handleDelete = (id: string) => {
		console.log('Delete', id);
	};

	//fetch first fetchResources
	useEffect(() => {
		dispatch(fetchResources({ size: rowsPerPage, searchAfter, searchBefore }));
	}, [dispatch, page, rowsPerPage]);

	const renderCustomContent = (setFilter: (filter: any) => void) => (
		<Box sx={{ mt: 1 }}>
			<TextField
				variant='outlined'
				label='Filter by Name'
				onChange={(e) => setFilter({ name: e.target.value })}
				fullWidth
				sx={{ mb: 1 }}
			/>
			<TextField
				variant='outlined'
				label='Filter by Code'
				onChange={(e) => setFilter({ code: e.target.value })}
				fullWidth
				sx={{ mb: 1 }}
			/>
		</Box>
	);

	return (
		<Box sx={{ flexGrow: 1 }}>
			<Grid2 container spacing={3}>
				{loading ? (
					<>
						<Skeleton variant='text' width='100%' height={120} />
						<Skeleton variant='rectangular' width='100%' height={40} />
						<Skeleton variant='rectangular' width='100%' height={40} />
						<Skeleton variant='rectangular' width='100%' height={40} />
						<Skeleton variant='rectangular' width='100%' height={40} />
						<Skeleton variant='rectangular' width='100%' height={40} />
						<Skeleton variant='text' width='100%' height={120} />
					</>
				) : (
					<TableComponent
						columns={columns}
						rows={resources?.items || []}
						onView={handleView}
						onEdit={handleEdit}
						onDelete={handleDelete}
						showFilter
						showSearch
						customContent={renderCustomContent}
						page={page}
						rowsPerPage={rowsPerPage}
						totalRows={resources?.total || 0}
						onPageChange={setPage}
						onRowsPerPageChange={setRowsPerPage}
					/>
				)}
			</Grid2>
		</Box>
	);
};

export default ResourcesPage;
