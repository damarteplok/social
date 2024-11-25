import { useState } from 'react';
import {
	Button,
	Modal,
	Box,
	Typography,
	IconButton,
	TextField,
} from '@mui/material';
import { PageContainerToolbar, useNotifications } from '@toolpad/core';
import { useFormik } from 'formik';
import { modalStyle } from '../../../layouts/modal';
import * as Yup from 'yup';
import { useDropzone } from 'react-dropzone';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';

const PageToolbarResource: React.FC = () => {
	const [open, setOpen] = useState(false);
	const [openProcessModal, setOpenProcessModal] = useState(false);
	const [variables, setVariables] = useState([{ key: '', value: '' }]);
	const notifications = useNotifications();

	const handleOpen = () => {
		formik.resetForm();
		processFormik.resetForm();
		setOpen(true);
	};
	const handleClose = () => {
		formik.resetForm();
		processFormik.resetForm();
		setOpen(false);
	};

	const handleOpenProcessModal = () => {
		setOpenProcessModal(true);
	};
	const handleCloseProcessModal = () => {
		setOpenProcessModal(false);
	};

	const formik = useFormik({
		enableReinitialize: true,
		initialValues: {
			bpmnFile: null as File | null,
			formFiles: [] as File[],
		},
		validationSchema: Yup.object({
			bpmnFile: Yup.mixed().required('Required'),
			formFiles: Yup.array(),
		}),
		onSubmit: async (values, { setSubmitting }) => {
			try {
				setSubmitting(true);
				const formData = new FormData();
				if (values.bpmnFile) {
					formData.append('files', values.bpmnFile);
				}
				values.formFiles.forEach((file) => {
					formData.append('files', file);
				});

				setSubmitting(false);
			} catch (err) {
				setSubmitting(false);
				notifications.show('Error uploading files', {
					severity: 'error',
					autoHideDuration: 3000,
				});
			}
		},
	});

	const processFormik = useFormik({
		initialValues: {
			processDefinitionKey: '',
			variables: [{ key: '', value: '' }],
		},
		validationSchema: Yup.object({
			processDefinitionKey: Yup.string().required('Required'),
		}),
		onSubmit: async (values, { setSubmitting }) => {
			// Handle process submission
			setSubmitting(false);
		},
	});

	const onDropBpmnFile = (acceptedFiles: File[]) => {
		if (acceptedFiles.length > 0) {
			formik.setFieldValue('bpmnFile', acceptedFiles[0]);
		} else {
			formik.setFieldError('bpmnFile', 'Invalid file type');
		}
	};

	const onDropFormFiles = (acceptedFiles: File[]) => {
		if (acceptedFiles.length > 0) {
			formik.setFieldValue('formFiles', [
				...formik.values.formFiles,
				...acceptedFiles,
			]);
		} else {
			formik.setFieldError('formFiles', 'Invalid file type');
		}
	};

	const handleDeleteFile = (
		fileToDelete: File,
		type: 'bpmn' | 'form',
		event: React.MouseEvent
	) => {
		event.stopPropagation();
		if (type === 'bpmn') {
			formik.setFieldValue('bpmnFile', null);
		} else {
			const updatedFiles = formik.values.formFiles.filter(
				(file) => file !== fileToDelete
			);
			formik.setFieldValue('formFiles', updatedFiles);
		}
	};

	const handleAddVariable = () => {
		setVariables([...variables, { key: '', value: '' }]);
		processFormik.setFieldValue('variables', [
			...variables,
			{ key: '', value: '' },
		]);
	};

	const handleRemoveVariable = (index: number) => {
		const updatedVariables = variables.filter((_, i) => i !== index);
		setVariables(updatedVariables);
		processFormik.setFieldValue('variables', updatedVariables);
	};

	const { getRootProps: getRootPropsBpmn, getInputProps: getInputPropsBpmn } =
		useDropzone({
			onDrop: onDropBpmnFile,
			accept: { 'application/xml': ['.bpmn', '.xml'] },
			maxFiles: 1,
		});

	const { getRootProps: getRootPropsForm, getInputProps: getInputPropsForm } =
		useDropzone({
			onDrop: onDropFormFiles,
			accept: { 'application/octet-stream': ['.form'] },
			multiple: true,
		});

	return (
		<PageContainerToolbar>
			<IconButton onClick={handleOpen}>
				<CloudUploadIcon />
			</IconButton>
			<IconButton onClick={handleOpenProcessModal}>
				<PlayArrowIcon />
			</IconButton>
			<Modal open={open} onClose={handleClose}>
				<Box sx={{ ...modalStyle }}>
					<Box
						sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}
						component={'form'}
						onSubmit={formik.handleSubmit}
					>
						<Box
							{...getRootPropsBpmn()}
							sx={{
								border: '1px dashed grey',
								padding: 2,
								textAlign: 'center',
							}}
						>
							<input {...getInputPropsBpmn()} />
							{formik.values.bpmnFile ? (
								<Box
									sx={{
										display: 'flex',
										alignItems: 'center',
										justifyContent: 'center',
									}}
								>
									<Typography>{formik.values.bpmnFile.name}</Typography>
									<Button
										onClick={(event) =>
											formik.values.bpmnFile &&
											handleDeleteFile(formik.values.bpmnFile, 'bpmn', event)
										}
									>
										Delete
									</Button>
								</Box>
							) : (
								<Box>
									Drag and drop a BPMN file here, or click to select one
								</Box>
							)}
							{formik.touched.bpmnFile && formik.errors.bpmnFile && (
								<Typography color='error'>{formik.errors.bpmnFile}</Typography>
							)}
						</Box>
						<Box
							{...getRootPropsForm()}
							sx={{
								border: '1px dashed grey',
								padding: 2,
								textAlign: 'center',
							}}
						>
							<input {...getInputPropsForm()} />
							{formik.values.formFiles.length > 0 ? (
								formik.values.formFiles.map((file, index) => (
									<Box
										key={index}
										sx={{
											display: 'flex',
											alignItems: 'center',
											justifyContent: 'center',
										}}
									>
										<Typography>{file.name}</Typography>
										<Button
											onClick={(event) => handleDeleteFile(file, 'form', event)}
										>
											Delete
										</Button>
									</Box>
								))
							) : (
								<Box>
									Drag and drop form files here, or click to select files
								</Box>
							)}
							{formik.touched.formFiles && formik.errors.formFiles && (
								<Typography color='error'>
									{Array.isArray(formik.errors.formFiles)
										? formik.errors.formFiles.join(', ')
										: formik.errors.formFiles}
								</Typography>
							)}
						</Box>
						<Button
							type='submit'
							variant='contained'
							disabled={formik.isSubmitting}
						>
							Deploy BPMN
						</Button>
					</Box>
				</Box>
			</Modal>
			<Modal open={openProcessModal} onClose={handleCloseProcessModal}>
				<Box sx={{ ...modalStyle, maxHeight: '80vh', overflowY: 'auto' }}>
					<Box
						sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}
						component={'form'}
						onSubmit={processFormik.handleSubmit}
					>
						<TextField
							label='Process Definition Key'
							name='processDefinitionKey'
							value={processFormik.values.processDefinitionKey}
							onChange={processFormik.handleChange}
							error={
								processFormik.touched.processDefinitionKey &&
								Boolean(processFormik.errors.processDefinitionKey)
							}
							helperText={
								processFormik.touched.processDefinitionKey &&
								processFormik.errors.processDefinitionKey
							}
						/>
						{variables.map((variable, index) => (
							<Box key={index} sx={{ display: 'flex', gap: 2 }}>
								<TextField
									label='Key'
									name={`variables[${index}].key`}
									value={variable.key}
									onChange={processFormik.handleChange}
								/>
								<TextField
									label='Value'
									name={`variables[${index}].value`}
									value={variable.value}
									onChange={processFormik.handleChange}
								/>
								<Button onClick={() => handleRemoveVariable(index)}>
									Delete
								</Button>
							</Box>
						))}
						<Button onClick={handleAddVariable}>Add Variable</Button>
						<Button
							type='submit'
							variant='contained'
							disabled={processFormik.isSubmitting}
						>
							Start Process
						</Button>
					</Box>
				</Box>
			</Modal>
		</PageContainerToolbar>
	);
};

export default PageToolbarResource;
