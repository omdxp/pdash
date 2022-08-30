export interface Order {
  id: string;
  supplier_id: string;
  customer_id: string;
  total_price: number;
  created_at: string;
  updated_at: string;
}

export type Orders = Order[];

export interface Supplier {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export type Suppliers = Supplier[];

export interface Customer {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export type Customers = Customer[];

export interface User {
  id: string;
  username: string;
  password: string;
  fullname: string;
  email: string;
  created_at: string;
  updated_at: string;
}
