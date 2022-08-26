import {
  ColumnDef,
  createSolidTable,
  flexRender,
  getCoreRowModel,
} from "@tanstack/solid-table";
import { Component, For, createSignal } from "solid-js";
import {
  deleteSupplier,
  refetchSuppliers,
  suppliers,
  updateSupplier,
} from "../../store";

import { Link } from "@solidjs/router";
import ModifySupplierModal from "../../components/modify-supplier-modal";
import { Supplier } from "../../interfaces";

const [selectedSupplier, setSelectedSupplier] = createSignal<Supplier>(null);

const defaultColumns: ColumnDef<Supplier>[] = [
  {
    accessorKey: "id",
    header: "ID",
    cell: (info) => (
      <Link href={`/suppliers/${info.getValue()}`}>
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
          setSelectedSupplier(
            suppliers().filter((el) => el.id === (info.getValue() as string))[0]
          );
        }}
      >
        Modify
      </button>
    ),
    footer: (info) => info.column.id,
  },
];

const Suppliers: Component = ({}) => {
  const table = createSolidTable({
    get data() {
      return suppliers();
    },
    columns: defaultColumns,
    getCoreRowModel: getCoreRowModel(),
  });
  return (
    <div class="flex flex-col">
      {suppliers.loading ? (
        <div>Loading...</div>
      ) : (
        <div>
          {selectedSupplier() && (
            <ModifySupplierModal
              supplier={selectedSupplier()}
              submit={async (draftCustomer) => {
                await updateSupplier(draftCustomer);
                await refetchSuppliers();
                setSelectedSupplier(null);
              }}
              cancel={() => setSelectedSupplier(null)}
              deleteSupplier={async () => {
                await deleteSupplier(selectedSupplier().id);
                await refetchSuppliers();
                setSelectedSupplier(null);
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
export default Suppliers;
