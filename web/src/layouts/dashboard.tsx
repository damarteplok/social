import { Outlet } from 'react-router-dom';
import { PageContainer } from '@toolpad/core/PageContainer';
import Typography from '@mui/material/Typography';
import {
	DashboardLayout,
	ThemeSwitcher,
	type SidebarFooterProps,
} from '@toolpad/core/DashboardLayout';
import { Tooltip, IconButton, Badge, Stack } from '@mui/material';
import NotificationsNoneIcon from '@mui/icons-material/NotificationsNone';

const ToolbarThemeSwitcherAndNotification = () => (
	<Stack direction='row'>
		<ThemeSwitcher />

		<Tooltip title='Notifications'>
			<IconButton color='inherit'>
				<Badge badgeContent={1} color='secondary'>
					<NotificationsNoneIcon />
				</Badge>
			</IconButton>
		</Tooltip>
	</Stack>
);

const SidebarFooter: React.FC<SidebarFooterProps> = ({ mini }) => {
	return (
		<Typography
			variant='caption'
			sx={{ m: 1, whiteSpace: 'nowrap', overflow: 'hidden' }}
		>
			{mini
				? '© Damarmunda'
				: `© ${new Date().getFullYear()} LowCode Camunda 8`}
		</Typography>
	);
};

export default function Layout() {
	return (
		<DashboardLayout
			slots={{
				toolbarActions: ToolbarThemeSwitcherAndNotification,
				sidebarFooter: SidebarFooter,
			}}
		>
			<PageContainer>
				<Outlet />
			</PageContainer>
		</DashboardLayout>
	);
}
