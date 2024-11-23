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
	Select,
	MenuItem,
	FormControl,
	InputLabel,
	Skeleton,
} from '@mui/material';
import {
	Search as SearchIcon,
	Edit as EditIcon,
	Delete as DeleteIcon,
	Visibility as ViewIcon,
	FilterList as FilterListIcon,
} from '@mui/icons-material';

interface TableComponentProps {
	columns: Array<{
		id: string;
		label: string;
		minWidth?: number;
		format?: (value: any) => string;
	}>;
	rows: Array<any>;
	onView?: (id: string) => void;
	onEdit?: (id: string) => void;
	onDelete?: (id: string) => void;
	onApplyFilter?: (filter: any) => void;
	onNextPage?: () => void;
	onPrevPage?: () => void;
	showSearch?: boolean;
	showFilter?: boolean;
	customContent?: React.ReactNode;
	page: number;
	rowsPerPage: number;
	totalRows: number;
	onPageChange: (newPage: number) => void;
	onRowsPerPageChange: (newRowsPerPage: number) => void;
	showActions?: boolean;
	showView?: boolean;
	showEdit?: boolean;
	showDelete?: boolean;
	customActions?: (row: any) => React.ReactNode;
	showNumber?: boolean;
	filter?: any;
	initFilter?: any;
	messageNoContent?: string;
	loading?: boolean;
}

const TableComponent: React.FC<TableComponentProps> = ({
	columns,
	rows,
	showSearch = false,
	showFilter = false,
	customContent,
	page,
	rowsPerPage,
	totalRows,
	showActions = true,
	showView = true,
	showEdit = true,
	showDelete = true,
	showNumber = true,
	filter,
	initFilter,
	messageNoContent = 'Data No Content',
	loading = false,
	onView = () => {},
	onEdit = () => {},
	onDelete = () => {},
	onApplyFilter = () => {},
	onNextPage = () => {},
	onPrevPage = () => {},
	onPageChange,
	onRowsPerPageChange,
	customActions,
}) => {
	const theme = useTheme();
	const [isFilterOpen, setIsFilterOpen] = useState(false);
	const [search, setSearch] = useState('');

	const handleReset = () => {
		setIsFilterOpen(false);
		onApplyFilter(initFilter);
	};

	const handleApply = () => {
		setIsFilterOpen(false);
		onApplyFilter(filter);
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

	return (
		<Paper sx={{ width: '100%', overflow: 'hidden' }}>
			<Box sx={{ p: 2, display: 'flex', justifyContent: 'flex-end' }}>
				{showSearch && (
					<TextField
						size='small'
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
			{loading ? (
				Array.from({ length: 10 }).map((_, index) => (
					<Skeleton
						variant='rectangular'
						width='100%'
						height={40}
						sx={{ mb: 1 }}
					/>
				))
			) : (
				<TableContainer sx={{ overflowX: 'auto' }}>
					<Table stickyHeader aria-label='sticky table' size='medium'>
						<TableHead>
							<TableRow>
								{showNumber && ( // Add this block
									<TableCell key={'number'} style={{ minWidth: 50 }}>
										No
									</TableCell>
								)}
								{columns.map((column) => (
									<TableCell
										key={column.id}
										style={{ minWidth: column.minWidth }}
									>
										{column.label}
									</TableCell>
								))}
								{showActions && (
									<TableCell
										key={'actions'}
										style={{
											position: 'sticky',
											right: 0,
											background: theme.palette.background.paper,
											minWidth: 170,
										}}
									>
										Actions
									</TableCell>
								)}
							</TableRow>
						</TableHead>
						<TableBody>
							{filteredRows.length === 0 ? (
								<TableRow>
									<TableCell
										colSpan={
											columns.length +
											(showActions ? 1 : 0) +
											(showNumber ? 1 : 0)
										}
										align='center'
									>
										{messageNoContent}
									</TableCell>
								</TableRow>
							) : (
								filteredRows.map((row, index) => (
									<TableRow hover tabIndex={-1} key={`${index}-${row.id}`}>
										{showNumber && ( // Add this block
											<TableCell key={`number-${index}`}>
												{page * rowsPerPage + index + 1}
											</TableCell>
										)}
										{columns.map((column) => {
											const value = row[column.id];
											return (
												<TableCell key={`${index}-${row.id}-${column.id}`}>
													{column.format ? column.format(value) : value}
												</TableCell>
											);
										})}
										{showActions && (
											<TableCell
												key={`${index}-${row.id}-actions`}
												style={{
													position: 'sticky',
													right: 0,
													background: theme.palette.background.paper,
												}}
											>
												{showView && (
													<IconButton onClick={() => onView(row.id)}>
														<ViewIcon />
													</IconButton>
												)}
												{showEdit && (
													<IconButton onClick={() => onEdit(row.id)}>
														<EditIcon />
													</IconButton>
												)}
												{showDelete && (
													<IconButton onClick={() => onDelete(row.id)}>
														<DeleteIcon />
													</IconButton>
												)}
												{customActions && customActions(row)}
											</TableCell>
										)}
									</TableRow>
								))
							)}
						</TableBody>
					</Table>
				</TableContainer>
			)}
			<Box sx={{ display: 'flex', justifyContent: 'space-between', p: 2 }}>
				<FormControl variant='outlined' size='small' sx={{ minWidth: 120 }}>
					<InputLabel>Rows per page</InputLabel>
					<Select
						value={rowsPerPage}
						onChange={(e) =>
							onRowsPerPageChange(parseInt(e.target.value as string, 10))
						}
						label='Rows per page'
					>
						<MenuItem value={10}>10</MenuItem>
						<MenuItem value={25}>25</MenuItem>
						<MenuItem value={50}>50</MenuItem>
						<MenuItem value={250}>250</MenuItem>
					</Select>
				</FormControl>
				<Box sx={{ display: 'flex', justifyContent: 'flex-end' }}>
					<IconButton
						onClick={() => {
							handleChangePage(null, page - 1);
							onPrevPage();
						}}
						disabled={page === 0}
					>
						{'<'}
					</IconButton>
					<IconButton
						onClick={() => {
							handleChangePage(null, page + 1);
							onNextPage();
						}}
						disabled={page >= Math.ceil(totalRows / rowsPerPage) - 1}
					>
						{'>'}
					</IconButton>
				</Box>
			</Box>
			<Drawer
				anchor='right'
				open={isFilterOpen}
				onClose={toggleFilterDrawer(false)}
			>
				<Box
					sx={{
						width: 350,
						p: 2,
						mt: 8,
						display: 'flex',
						flexDirection: 'column',
						height: '100%',
					}}
				>
					<Typography variant='h6'>Filter</Typography>
					{/* Add your filter options here */}
					{customContent}
					<Box sx={{ flexGrow: 1 }} />
					<Box sx={{ display: 'flex', justifyContent: 'space-between', p: 2 }}>
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
