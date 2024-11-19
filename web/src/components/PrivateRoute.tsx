import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { checkAuth } from '../slices/modules/session/thunk';
import { RootState } from '../slices/store/rootReducer';

const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({
	children,
}) => {
	// TODO: Implement PrivateRoute component
	// const navigate = useNavigate();
	// const dispatch = useDispatch<any>();
	// const { session } = useSelector(
	// 	(state: RootState) => state.session
	// );

	// useEffect(() => {
	// 	const authenticate = async () => {
	// 		try {
	// 			await dispatch(checkAuth()).unwrap();
	// 		} catch {
	// 			navigate('/sign-in');
	// 		}
	// 	};

	// 	authenticate();
	// }, [dispatch, navigate]);

	return <>{children}</>;
};

export default PrivateRoute;
