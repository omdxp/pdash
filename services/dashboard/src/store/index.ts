import { CUSTOMERS_URL, ORDERS_URL, SUPPLIERS_URL } from "../config";
import { Customers, Orders, Suppliers } from "../interfaces";

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
