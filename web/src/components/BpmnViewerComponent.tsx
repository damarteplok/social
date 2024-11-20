import React, { useEffect, useRef } from 'react';
import BpmnModeler from 'bpmn-js/lib/NavigatedViewer';

interface BpmnViewerProps {
	xml: string;
	taskCounts: Record<string, number>; // Task counts keyed by task ID
}

const BpmnViewerComponent: React.FC<BpmnViewerProps> = ({
	xml,
	taskCounts,
}) => {
	const viewerRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		if (!viewerRef.current) return;

		const modeler = new BpmnModeler({
			container: viewerRef.current,
		});

		const importDiagram = async () => {
			try {
				const result = await modeler.importXML(xml);
				const { warnings } = result;

				const canvas = modeler.get('canvas') as any;
				const overlays = modeler.get('overlays') as any;
				const elementRegistry = modeler.get('elementRegistry') as any;

				// Add task counts to elements
				Object.keys(taskCounts).forEach((taskId) => {
					const element = elementRegistry.get(taskId);
					if (element) {
						const overlayHtml = document.createElement('div');
						overlayHtml.className = 'task-count-overlay';
						overlayHtml.innerText = `Active: ${taskCounts[taskId]}`;
						overlays.add(taskId, {
							position: {
								bottom: 0,
								right: 0,
							},
							html: overlayHtml,
						});
					}
				});

				canvas.zoom('fit-viewport');
			} catch (err) {
				console.log('something went wrong:');
			}
		};

		importDiagram();

		return () => {
			modeler.destroy();
		};
	}, [xml, taskCounts]);

	return <div ref={viewerRef} style={{ width: '100%', height: '600px' }} />;
};

export default BpmnViewerComponent;
