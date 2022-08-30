import { AUTH_URL, CUSTOMERS_URL, ORDERS_URL, SUPPLIERS_URL } from "../config";
import {
  Customer,
  Customers,
  Order,
  Orders,
  Supplier,
  Suppliers,
  User,
} from "../interfaces";
import { Signal, createResource, createSignal } from "solid-js";

export const [token, setToken] = createStoredSignal("token", null);

async function fetchCustomers(): Promise<Customers> {
  if (token()) {
    const res = await fetch(CUSTOMERS_URL, {
      headers: {
        Authorization: `Bearer ${token()}`,
      },
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
    return [];
  } else {
    return [];
  }
}

async function fetchOrders(): Promise<Orders> {
  if (token()) {
    const res = await fetch(ORDERS_URL, {
      headers: {
        Authorization: `Bearer ${token()}`,
      },
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
    return [];
  } else {
    return [];
  }
}

async function fetchSuppliers(): Promise<Suppliers> {
  if (token()) {
    const res = await fetch(SUPPLIERS_URL, {
      headers: {
        Authorization: `Bearer ${token()}`,
      },
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
    return [];
  } else {
    return [];
  }
}

export const [customers, { refetch: refetchCustomers }] =
  createResource(fetchCustomers);
export const [orders, { refetch: refetchOrders }] = createResource(fetchOrders);
export const [suppliers, { refetch: refetchSuppliers }] =
  createResource(fetchSuppliers);

export async function fetchSupplier(id: string): Promise<Supplier> {
  if (token()) {
    const res = await fetch(SUPPLIERS_URL + "/" + id, {
      headers: {
        Authorization: `Bearer ${token()}`,
      },
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
    return null;
  } else {
    return null;
  }
}

export async function fetchCustomer(id: string): Promise<Customer> {
  if (token()) {
    const res = await fetch(CUSTOMERS_URL + "/" + id, {
      headers: {
        Authorization: `Bearer ${token()}`,
      },
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
    return null;
  } else {
    return null;
  }
}

export async function fetchSupplierOrders(id: string): Promise<Orders> {
  if (token()) {
    const res = await fetch(ORDERS_URL + "?supplier_id=" + id, {
      headers: {
        Authorization: `Bearer ${token()}`,
      },
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
    return [];
  } else {
    return [];
  }
}

export async function fetchCustomerOrders(id: string): Promise<Orders> {
  if (token()) {
    const res = await fetch(ORDERS_URL + "?customer_id=" + id, {
      headers: {
        Authorization: `Bearer ${token()}`,
      },
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
    return [];
  } else {
    return [];
  }
}

export async function updateOrder(order: Order) {
  if (token()) {
    const res = await fetch(ORDERS_URL + "/" + order.id, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token()}`,
      },
      body: JSON.stringify(order),
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
  }
}

export async function deleteOrder(id: string) {
  if (token()) {
    const res = await fetch(ORDERS_URL + "/" + id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token()}`,
      },
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
  }
}

export async function updateCustomer(customer: Customer) {
  if (token()) {
    const res = await fetch(CUSTOMERS_URL + "/" + customer.id, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token()}`,
      },
      body: JSON.stringify(customer),
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
  }
}

export async function deleteCustomer(id: string) {
  if (token()) {
    const res = await fetch(CUSTOMERS_URL + "/" + id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token()}`,
      },
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
  }
}

export async function updateSupplier(supplier: Supplier) {
  if (token()) {
    const res = await fetch(SUPPLIERS_URL + "/" + supplier.id, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token()}`,
      },
      body: JSON.stringify(supplier),
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
  }
}

export async function deleteSupplier(id: string) {
  if (token()) {
    const res = await fetch(SUPPLIERS_URL + "/" + id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token()}`,
      },
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
  }
}

export async function addCustomer(customer: Partial<Customer>) {
  if (token()) {
    const res = await fetch(CUSTOMERS_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token()}`,
      },
      body: JSON.stringify(customer),
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
  }
}

export async function addSupplier(supplier: Partial<Supplier>) {
  if (token()) {
    const res = await fetch(SUPPLIERS_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token()}`,
      },
      body: JSON.stringify(supplier),
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
  }
}

export async function addOrder(order: Partial<Order>) {
  if (token()) {
    const res = await fetch(ORDERS_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token()}`,
      },
      body: JSON.stringify(order),
    });
    if (res.status == 200) {
      return await res.json();
    }
    if (res.status == 401 || res.status == 500) {
      setToken(null);
    }
  }
}

export async function login(user: Partial<User>) {
  const res = await fetch(`${AUTH_URL}/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(user),
  });
  if (res.status == 200) {
    return await res.json();
  }
}

export async function signup(user: Partial<User>) {
  const res = await fetch(AUTH_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(user),
  });
  if (res.status == 201) {
    return await res.json();
  }
}

function createStoredSignal<T>(
  key: string,
  defaultValue: T,
  storage = localStorage
): Signal<T> {
  const initialValue = storage.getItem(key)
    ? (JSON.parse(storage.getItem(key)) as T)
    : defaultValue;

  const [value, setValue] = createSignal<T>(initialValue);

  const setValueAndStore = ((arg) => {
    const v = setValue(arg);
    storage.setItem(key, JSON.stringify(v));
    return v;
  }) as typeof setValue;

  return [value, setValueAndStore];
}
