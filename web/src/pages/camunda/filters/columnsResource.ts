export const columnsResources = [
	{ id: 'bpmnProcessId', label: 'Bpmn Process Id', minWidth: 170 },
	{ id: 'key', label: 'Key', minWidth: 170 },
	{
		id: 'processDefinitionKey',
		label: 'Process Definition Key',
		minWidth: 200,
	},
	{ id: 'processVersion', label: 'Version' },
	{
		id: 'startDate',
		label: 'Start Date',
		format: (value: string) => new Date(value).toLocaleDateString('en-ID'),
	},
	{
		id: 'endDate',
		label: 'End Date',
		format: (value: string) => new Date(value).toLocaleDateString('en-ID'),
	},
	{ id: 'state', label: 'State', minWidth: 170 },
	{
		id: 'incident',
		label: 'Incident',
		format: (value: any) => {
			return value ? 'TRUE' : 'FALSE';
		},
	},
];

export interface RowResource {
	bpmnProcessId: string;
	endDate: string;
	incident: boolean;
	key: number;
	processDefinitionKey: number;
	processVersion: number;
	startDate: string;
	state: string;
	tenantId: string;
}
