import {
  ColumnDef,
  createSolidTable,
  flexRender,
  getCoreRowModel,
} from "@tanstack/solid-table";
import { Component, For, createSignal } from "solid-js";
import { deleteOrder, orders, refetchOrders, updateOrder } from "../../store";

import DeleteModal from "../../components/delete-modal";
import { Link } from "@solidjs/router";
import ModifyOrderModal from "../../components/modify-order-modal";
import { Order } from "../../interfaces";

const [selectedOrder, setSelectedOrder] = createSignal<Order>(null);

const defaultColumns: ColumnDef<Order>[] = [
  {
    accessorKey: "id",
    header: "ID",
    cell: (info) => (
      <button
        onClick={() =>
          setSelectedOrder(
            orders().filter((el) => el.id === (info.getValue() as string))[0]
          )
        }
      >
        {info.getValue() as string}
      </button>
    ),
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

const Home: Component = ({}) => {
  const table = createSolidTable({
    get data() {
      return orders();
    },
    columns: defaultColumns,
    getCoreRowModel: getCoreRowModel(),
  });
  return (
    <div class="flex flex-col">
      {orders.loading ? (
        <div>Loading...</div>
      ) : (
        <div>
          {selectedOrder() && (
            <ModifyOrderModal
              order={selectedOrder()}
              submit={async (draftOrder) => {
                await updateOrder(draftOrder);
                await refetchOrders();
                setSelectedOrder(null);
              }}
              cancel={() => setSelectedOrder(null)}
              deleteOrder={async () => {
                await deleteOrder(selectedOrder().id);
                await refetchOrders();
                setSelectedOrder(null);
              }}
            />
          )}
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
        </div>
      )}
    </div>
  );
};
export default Home;
