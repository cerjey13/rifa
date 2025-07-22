import { type FC } from 'react';
import {
  List,
  Datagrid,
  TextField,
  NumberField,
  DateField,
  FunctionField,
  type ListProps,
} from 'react-admin';

export const PurchaseList: FC<ListProps> = (props) => (
  <List {...props}>
    <Datagrid rowClick='show'>
      <TextField source='user_id' label='Usuario' />
      <NumberField source='quantity' label='Cantidad' />
      <NumberField source='monto_bs' label='Monto (Bs)' />
      <NumberField source='monto_usd' label='Monto (USD)' />
      <TextField source='payment_method' label='Método de pago' />
      <TextField source='transaction_digits' label='Transacción' />
      <DateField source='created_at' label='Fecha' showTime />
      <TextField source='status' label='Estado' />
      <FunctionField
        label='Comprobante'
        render={(record: Purchase) => (
          <img
            src={`data:image/png;base64,${record.paymentScreenshot}`}
            alt='Comprobante'
            style={{ maxHeight: 60, maxWidth: 60, objectFit: 'contain' }}
          />
        )}
      />
    </Datagrid>
  </List>
);
