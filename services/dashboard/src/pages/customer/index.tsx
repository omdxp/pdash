import {
  ColumnDef,
  createSolidTable,
  flexRender,
  getCoreRowModel,
} from "@tanstack/solid-table";
import { Component, For, createResource } from "solid-js";
import { Link, useParams } from "@solidjs/router";
import { fetchCustomer, fetchCustomerOrders } from "../../store";

import { Order } from "../../interfaces";

const defaultColumns: ColumnDef<Order>[] = [
  {
    accessorKey: "id",
    header: "ID",
    cell: (info) => info.getValue(),
    footer: (info) => info.column.id,
  },
  {
    accessorKey: "total_price",
    header: "Total price",
    cell: (info) => info.getValue(),
    footer: (info) => info.column.id,
  },
  {
    accessorKey: "customer_id",
    header: "Customer ID",
    cell: (info) => (
      <Link href={`/customers/${info.getValue()}`}>
        <a class="hover:text-blue-600">{info.getValue() as string}</a>
      </Link>
    ),
    footer: (info) => info.column.id,
  },
  {
    accessorKey: "supplier_id",
    header: "Supplier ID",
    cell: (info) => (
      <Link href={`/suppliers/${info.getValue()}`}>
        <a class="hover:text-blue-600">{info.getValue() as string}</a>
      </Link>
    ),
    footer: (info) => info.column.id,
  },
  {
    accessorKey: "created_at",
    header: "Created at",
    cell: (info) => (
      <span>
        {new Date(info.getValue() as string).toDateString()} -{" "}
        {new Date(info.getValue() as string).toTimeString()}
      </span>
    ),
    footer: (info) => info.column.id,
  },
  {
    accessorKey: "updated_at",
    header: "Updated at",
    cell: (info) => (
      <span>
        {new Date(info.getValue() as string).toDateString()} -{" "}
        {new Date(info.getValue() as string).toTimeString()}
      </span>
    ),
    footer: (info) => info.column.id,
  },
];

const Customer: Component = ({}) => {
  const params = useParams();
  const [customerData] = createResource(() => params.id, fetchCustomer);
  const [customerOrders] = createResource(() => params.id, fetchCustomerOrders);
  const table = createSolidTable({
    get data() {
      return customerOrders();
    },
    columns: defaultColumns,
    getCoreRowModel: getCoreRowModel(),
  });
  return (
    <div>
      {customerData() && (
        <>
          <p>ID: {customerData().id}</p>
          <p>Name: {customerData().name}</p>
          <p>
            Created at: {new Date(customerData().created_at).toDateString()} -{" "}
            {new Date(customerData().created_at).toTimeString()}
          </p>
          <p>
            Updated at: {new Date(customerData().updated_at).toDateString()} -{" "}
            {new Date(customerData().updated_at).toTimeString()}
          </p>
        </>
      )}
      <h2>Orders</h2>
      {customerOrders() && (
        <table>
          <thead>
            <For each={table.getHeaderGroups()}>
              {(headerGroup) => (
                <tr>
                  <For each={headerGroup.headers}>
                    {(header) => (
                      <th>
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                              header.column.columnDef.header,
                              header.getContext()
                            )}
                      </th>
                    )}
                  </For>
                </tr>
              )}
            </For>
          </thead>
          <tbody>
            <For each={table.getRowModel().rows}>
              {(row) => (
                <tr>
                  <For each={row.getVisibleCells()}>
                    {(cell) => (
                      <td>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext()
                        )}
                      </td>
                    )}
                  </For>
                </tr>
              )}
            </For>
          </tbody>
          <tfoot>
            <For each={table.getFooterGroups()}>
              {(footerGroup) => (
                <tr>
                  <For each={footerGroup.headers}>
                    {(header) => (
                      <th>
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                              header.column.columnDef.footer,
                              header.getContext()
                            )}
                      </th>
                    )}
                  </For>
                </tr>
              )}
            </For>
          </tfoot>
        </table>
      )}
    </div>
  );
};
export default Customer;
