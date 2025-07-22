import { Admin, Layout, Menu, Resource, type LayoutProps } from 'react-admin';
import { PurchaseList } from '@src/admin/components/Purchases';
import { Dashboard } from './components/Dashboard';
import { FaShoppingCart, FaHome } from 'react-icons/fa';
import { AiOutlineLogout } from 'react-icons/ai';
import { useAuth } from '@src/context/useAuth';

const CustomMenu = () => {
  const { logout } = useAuth();
  return (
    <Menu>
      <Menu.DashboardItem
        tabIndex={1}
        to='/dashboard'
        primaryText='Compras'
        leftIcon={<FaShoppingCart />}
      />
      <Menu.DashboardItem
        tabIndex={2}
        to='/'
        primaryText='Pagina Principal'
        leftIcon={<FaHome />}
      />
      <Menu.DashboardItem
        tabIndex={3}
        onClick={logout}
        primaryText='Cerrar Sesion'
        leftIcon={<AiOutlineLogout />}
      />
    </Menu>
  );
};

const CustomLayout = (props: LayoutProps) => (
  <Layout {...props} menu={CustomMenu} className='Menu' />
);

const AdminApp = () => (
  <Admin dashboard={Dashboard} layout={CustomLayout}>
    <Resource name='dashboard' list={PurchaseList} />
  </Admin>
);
export default AdminApp;
