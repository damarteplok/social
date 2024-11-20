import { useState } from 'react';
import { Button, Modal, Box, TextField } from '@mui/material';
import SchemaIcon from '@mui/icons-material/Schema';
import FileDownloadIcon from '@mui/icons-material/FileDownload';
import { PageContainerToolbar, useNotifications } from '@toolpad/core';
import { useFormik } from 'formik';
import axiosInstance from '../../../utils/axiosInstance';
import { modalStyle } from '../../../layouts/modal';
import BpmnViewerComponent from '../../../components/BpmnViewerComponent';
import * as Yup from 'yup';

const PageToolbarResource: React.FC = () => {
	const [open, setOpen] = useState(false);
	const [bpmnXML, setBpmnXML] = useState('');
	const [taskCounts, setTaskCounts] = useState<Record<string, number>>({});
	const notifications = useNotifications();

	const handleOpen = () => {
		setBpmnXML('');
		formik.resetForm();
		setOpen(true);
	};
	const handleClose = () => {
		setBpmnXML('');
		formik.resetForm();
		setOpen(false);
	};

	const formik = useFormik({
		enableReinitialize: true,
		initialValues: {
			processDefinitionKey: '',
		},
		validationSchema: Yup.object({
			processDefinitionKey: Yup.string().required('Required'),
		}),
		onSubmit: async (values, { setSubmitting }) => {
			try {
				setSubmitting(true);
				const response = await axiosInstance.get(
					`/camunda/resource/${values.processDefinitionKey}/xml`
				);
				const xml = response.data;
				setBpmnXML(xml);
				setSubmitting(false);
			} catch (err) {
				setSubmitting(false);
				notifications.show('Error fetching BPMN', {
					severity: 'error',
					autoHideDuration: 3000,
				});
			}
		},
	});

	return (
		<PageContainerToolbar>
			<Button startIcon={<SchemaIcon />} color='inherit' onClick={handleOpen}>
				BPMN
			</Button>
			<Button startIcon={<FileDownloadIcon />} color='inherit'>
				Export
			</Button>
			<Modal open={open} onClose={handleClose}>
				<Box sx={{ ...modalStyle }}>
					{bpmnXML && (
						<BpmnViewerComponent xml={bpmnXML} taskCounts={taskCounts} />
					)}

					<Box
						sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}
						component={'form'}
						onSubmit={formik.handleSubmit}
					>
						<TextField
							required
							label='Process Definition Key'
							name='processDefinitionKey'
							value={formik.values.processDefinitionKey}
							onChange={formik.handleChange}
							onBlur={formik.handleBlur}
							size='small'
							error={
								formik.touched.processDefinitionKey &&
								Boolean(formik.errors.processDefinitionKey)
							}
							helperText={
								formik.touched.processDefinitionKey &&
								formik.errors.processDefinitionKey
							}
							fullWidth
						/>
						<Button
							type='submit'
							variant='contained'
							disabled={formik.isSubmitting}
						>
							Fetch BPMN
						</Button>
					</Box>
				</Box>
			</Modal>
		</PageContainerToolbar>
	);
};

export default PageToolbarResource;
