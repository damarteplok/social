import React from 'react';
import {
	Box,
	TextField,
	FormControl,
	InputLabel,
	MenuItem,
	Select,
} from '@mui/material';

interface FilterResourceComponentProps {
	filter: any;
	setFilter: (filter: any) => void;
}

const FilterResourceComponent: React.FC<FilterResourceComponentProps> = ({
	filter,
	setFilter,
}) => {
	return (
		<Box sx={{ mt: 1 }}>
			<FormControl fullWidth sx={{ mb: 2 }}>
				<TextField
					variant='outlined'
					label='Filter by Bpmn Process Id'
					onChange={(e) =>
						setFilter((prev: any) => ({
							...prev,
							bpmnProcessId: e.target.value,
						}))
					}
					value={filter.bpmnProcessId}
					fullWidth
					size='small'
				/>
			</FormControl>
			<FormControl fullWidth sx={{ mb: 2 }}>
				<TextField
					variant='outlined'
					label='Filter by Start Date'
					onChange={(e) =>
						setFilter((prev: any) => ({ ...prev, startDate: e.target.value }))
					}
					fullWidth
					value={filter.startDate}
					size='small'
				/>
			</FormControl>
			<FormControl fullWidth sx={{ mb: 2 }}>
				<TextField
					variant='outlined'
					label='Filter by End Date'
					onChange={(e) =>
						setFilter((prev: any) => ({ ...prev, endDate: e.target.value }))
					}
					fullWidth
					value={filter.endDate}
					size='small'
				/>
			</FormControl>
			<FormControl fullWidth sx={{ mb: 2 }}>
				<TextField
					variant='outlined'
					label='Filter by Process Definition Key'
					onChange={(e) =>
						setFilter((prev: any) => ({
							...prev,
							processDefinitionKey: e.target.value,
						}))
					}
					fullWidth
					value={filter.processDefinitionKey}
					size='small'
				/>
			</FormControl>
			<FormControl fullWidth sx={{ mb: 2 }}>
				<TextField
					variant='outlined'
					label='Filter by Parent Process Instance Key'
					onChange={(e) =>
						setFilter((prev: any) => ({
							...prev,
							parentProcessInstanceKey: e.target.value,
						}))
					}
					fullWidth
					value={filter.parentProcessInstanceKey}
					size='small'
				/>
			</FormControl>
			<FormControl fullWidth sx={{ mb: 2 }}>
				<InputLabel size='small'>Filter by State</InputLabel>
				<Select
					variant='outlined'
					label='Filter by State'
					onChange={(e) =>
						setFilter((prev: any) => ({ ...prev, state: e.target.value }))
					}
					fullWidth
					defaultValue=''
					size='small'
					value={filter.state}
				>
					<MenuItem value=''>None</MenuItem>
					<MenuItem value='ACTIVE'>Active</MenuItem>
					<MenuItem value='COMPLETED'>Completed</MenuItem>
					<MenuItem value='CANCELED'>Canceled</MenuItem>
				</Select>
			</FormControl>
		</Box>
	);
};

export default FilterResourceComponent;
