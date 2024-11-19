import React, { useState } from 'react';
import {
	Box,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Paper,
	IconButton,
	TextField,
	InputAdornment,
	Drawer,
	Button,
	Typography,
	useTheme,
} from '@mui/material';
import {
	Search as SearchIcon,
	Edit as EditIcon,
	Delete as DeleteIcon,
	Visibility as ViewIcon,
	FilterList as FilterListIcon,
} from '@mui/icons-material';

interface TableComponentProps {
	columns: Array<{ id: string; label: string; minWidth?: number }>;
	rows: Array<any>;
	onView: (id: string) => void;
	onEdit: (id: string) => void;
	onDelete: (id: string) => void;
	showSearch?: boolean;
	showFilter?: boolean;
	customContent?: (setFilter: (filter: any) => void) => React.ReactNode;
	page: number;
	rowsPerPage: number;
	totalRows: number;
	onPageChange: (newPage: number) => void;
	onRowsPerPageChange: (newRowsPerPage: number) => void;
}

const TableComponent: React.FC<TableComponentProps> = ({
	columns,
	rows,
	onView,
	onEdit,
	onDelete,
	showSearch = false,
	showFilter = false,
	customContent,
	page,
	rowsPerPage,
	totalRows,
	onPageChange,
	onRowsPerPageChange,
}) => {
	const [search, setSearch] = useState('');
	const theme = useTheme();
	const [filter, setFilter] = useState({});
	const [isFilterOpen, setIsFilterOpen] = useState(false);

	const handleReset = () => {
		setFilter({});
		setIsFilterOpen(false);
	};

	const handleApply = () => {
		setIsFilterOpen(false);
	};

	const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		setSearch(event.target.value);
	};

	const toggleFilterDrawer = (open: boolean) => () => {
		setIsFilterOpen(open);
	};

	const filteredRows = rows.filter((row) =>
		columns.some((column) =>
			row[column.id].toString().toLowerCase().includes(search.toLowerCase())
		)
	);

	const handleChangePage = (event: unknown, newPage: number) => {
		onPageChange(newPage);
	};

	const handleChangeRowsPerPage = (
		event: React.ChangeEvent<HTMLInputElement>
	) => {
		onRowsPerPageChange(parseInt(event.target.value, 10));
	};

	return (
		<Paper sx={{ width: '100%', overflow: 'hidden' }}>
			<Box sx={{ p: 2, display: 'flex', justifyContent: 'flex-end' }}>
				{showSearch && (
					<TextField
						sx={{
							mr: 1,
						}}
						variant='outlined'
						placeholder='Search...'
						value={search}
						onChange={handleSearchChange}
						InputProps={{
							startAdornment: (
								<InputAdornment position='start'>
									<SearchIcon />
								</InputAdornment>
							),
						}}
					/>
				)}
				{showFilter && (
					<Button
						sx={{
							mr: 1,
						}}
						variant='outlined'
						startIcon={<FilterListIcon />}
						onClick={toggleFilterDrawer(true)}
					>
						Filter
					</Button>
				)}
			</Box>
			<TableContainer sx={{ overflowX: 'auto' }}>
				<Table stickyHeader aria-label='sticky table' size='small'>
					<TableHead>
						<TableRow>
							{columns.map((column) => (
								<TableCell
									key={column.id}
									style={{ minWidth: column.minWidth }}
								>
									{column.label}
								</TableCell>
							))}
							<TableCell
								style={{
									position: 'sticky',
									right: 0,
									background: theme.palette.background.paper,
									minWidth: 170,
								}}
							>
								Actions
							</TableCell>
						</TableRow>
					</TableHead>
					<TableBody>
						{filteredRows
							.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
							.map((row) => (
								<TableRow hover tabIndex={-1} key={row.id}>
									{columns.map((column) => {
										const value = row[column.id];
										return <TableCell key={column.id}>{value}</TableCell>;
									})}
									<TableCell
										style={{
											position: 'sticky',
											right: 0,
											background: theme.palette.background.paper,
										}}
									>
										<IconButton onClick={() => onView(row.id)}>
											<ViewIcon />
										</IconButton>
										<IconButton onClick={() => onEdit(row.id)}>
											<EditIcon />
										</IconButton>
										<IconButton onClick={() => onDelete(row.id)}>
											<DeleteIcon />
										</IconButton>
									</TableCell>
								</TableRow>
							))}
					</TableBody>
				</Table>
			</TableContainer>
			<Box sx={{ display: 'flex', justifyContent: 'flex-end', p: 2 }}>
				<IconButton
					onClick={() => handleChangePage(null, page - 1)}
					disabled={page === 0}
				>
					{'<'}
				</IconButton>
				<IconButton
					onClick={() => handleChangePage(null, page + 1)}
					disabled={page >= Math.ceil(totalRows / rowsPerPage) - 1}
				>
					{'>'}
				</IconButton>
			</Box>
			<Drawer
				anchor='right'
				open={isFilterOpen}
				onClose={toggleFilterDrawer(false)}
			>
				<Box sx={{ width: 350, p: 2, mt: 8 }}>
					<Typography variant='h6'>Filter</Typography>
					{/* Add your filter options here */}
					{customContent && customContent(setFilter)}
					<Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
						<Button variant='text' onClick={handleReset} fullWidth>
							Reset
						</Button>
						<Button variant='contained' onClick={handleApply} fullWidth>
							Apply
						</Button>
					</Box>
				</Box>
			</Drawer>
		</Paper>
	);
};

export default TableComponent;
