import React from 'react';
import { Fab } from '@mui/material';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import ScrollToTop from 'react-scroll-to-top';

const ScrollToTopButton: React.FC = () => {
	return (
		<ScrollToTop
			smooth
            width='150px'
            height='150px'
            top={100}
			// component={
			// 	<Fab size='small' aria-label='scroll back to top'>
			// 		<KeyboardArrowUpIcon />
			// 	</Fab>
			// }
		/>
	);
};

export default ScrollToTopButton;
