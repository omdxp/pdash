import {
  ColumnDef,
  createSolidTable,
  flexRender,
  getCoreRowModel,
} from "@tanstack/solid-table";
import { Component, For, createSignal, onMount } from "solid-js";
import {
  addCustomer,
  customers,
  deleteCustomer,
  refetchCustomers,
  updateCustomer,
} from "../../store";

import AddCustomerModal from "../../components/add-customer-modal";
import { Customer } from "../../interfaces";
import { Link } from "@solidjs/router";
import ModifyCustomerModal from "../../components/modify-customer-modal";

const [selectedCustomer, setSelectedCustomer] = createSignal<Customer>(null);

const defaultColumns: ColumnDef<Customer>[] = [
  {
    accessorKey: "id",
    header: "ID",
    cell: (info) => (
      <Link href={`/customers/${info.getValue()}`}>
        <a class="hover:text-blue-600">{info.getValue() as string}</a>
      </Link>
    ),
    footer: (info) => info.column.id,
  },
  {
    accessorKey: "name",
    header: "Name",
    cell: (info) => info.getValue(),
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
  {
    accessorKey: "id",
    header: "Modify",
    cell: (info) => (
      <button
        class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        onClick={() => {
          setSelectedCustomer(
            customers().filter((el) => el.id === (info.getValue() as string))[0]
          );
        }}
      >
        Modify
      </button>
    ),
    footer: (info) => info.column.id,
  },
];

const Customers: Component = ({}) => {
  const table = createSolidTable({
    get data() {
      return customers();
    },
    columns: defaultColumns,
    getCoreRowModel: getCoreRowModel(),
  });
  const [addCustomerModalShown, setAddCustomerModalShown] = createSignal(false);
  onMount(() => {
    refetchCustomers();
  });
  return (
    <div class="flex flex-col">
      {customers.loading ? (
        <div>Loading...</div>
      ) : (
        <div>
          <button
            class="bg-blue-500 m-2 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
            onClick={() => setAddCustomerModalShown(true)}
          >
            Add new customer
          </button>
          {addCustomerModalShown() && (
            <AddCustomerModal
              submit={async (draftCustomer) => {
                await addCustomer(draftCustomer);
                await refetchCustomers();
                setAddCustomerModalShown(false);
              }}
              cancel={() => setAddCustomerModalShown(false)}
            />
          )}
          {selectedCustomer() && (
            <ModifyCustomerModal
              customer={selectedCustomer()}
              submit={async (draftCustomer) => {
                await updateCustomer(draftCustomer);
                await refetchCustomers();
                setSelectedCustomer(null);
              }}
              cancel={() => setSelectedCustomer(null)}
              deleteCustomer={async () => {
                await deleteCustomer(selectedCustomer().id);
                await refetchCustomers();
                setSelectedCustomer(null);
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
export default Customers;
