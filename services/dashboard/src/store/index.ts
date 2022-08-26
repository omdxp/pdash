import { CUSTOMERS_URL, ORDERS_URL, SUPPLIERS_URL } from "../config";
import {
  Customer,
  Customers,
  Orders,
  Supplier,
  Suppliers,
} from "../interfaces";

import { createResource } from "solid-js";

async function fetchCustomers(): Promise<Customers> {
  const res = await fetch(CUSTOMERS_URL);
  if (res.status == 200) {
    return await res.json();
  }
  return [];
}

async function fetchOrders(): Promise<Orders> {
  const res = await fetch(ORDERS_URL);
  if (res.status == 200) {
    return await res.json();
  }
  return [];
}

async function fetchSuppliers(): Promise<Suppliers> {
  const res = await fetch(SUPPLIERS_URL);
  if (res.status == 200) {
    return await res.json();
  }
  return [];
}

export const [customers, { mutate, refetch }] = createResource(fetchCustomers);
export const [orders, { mutate: mutateOrders }] = createResource(fetchOrders);
export const [suppliers, { mutate: mutateSuppliers }] =
  createResource(fetchSuppliers);

export async function fetchSupplier(id: string): Promise<Supplier> {
  const res = await fetch(SUPPLIERS_URL + "/" + id);
  if (res.status == 200) {
    return await res.json();
  }
  return null;
}

export async function fetchCustomer(id: string): Promise<Customer> {
  const res = await fetch(CUSTOMERS_URL + "/" + id);
  if (res.status == 200) {
    return await res.json();
  }
  return null;
}

export async function fetchSupplierOrders(id: string): Promise<Orders> {
  const res = await fetch(ORDERS_URL + "?supplier_id=" + id);
  if (res.status == 200) {
    return await res.json();
  }
  return [];
}

export async function fetchCustomerOrders(id: string): Promise<Orders> {
  const res = await fetch(ORDERS_URL + "?customer_id=" + id);
  if (res.status == 200) {
    return await res.json();
  }
  return [];
}
