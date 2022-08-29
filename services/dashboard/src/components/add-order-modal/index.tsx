import { Component, For, createSignal } from "solid-js";
import { customers, suppliers } from "../../store";

import { Order } from "../../interfaces";

interface AddOrderModalProps {
  submit: (draftOrder: Partial<Order>) => void;
  cancel: () => void;
}

const AddOrderModal: Component<AddOrderModalProps> = ({ submit, cancel }) => {
  const [draftOrder, setDraftOrder] = createSignal<Partial<Order>>({
    total_price: 0,
    customer_id: customers()[0].id,
    supplier_id: suppliers()[0].id,
  });
  return (
    <div
      id="authentication-modal"
      tabindex="-1"
      aria-hidden="true"
      class="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 w-full md:inset-0 h-modal md:h-full"
    >
      <div class="relative p-4 w-full max-w-md h-full md:h-auto">
        <div class="relative bg-white rounded-lg shadow dark:bg-gray-700">
          <button
            type="button"
            class="absolute top-3 right-2.5 text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center dark:hover:bg-gray-800 dark:hover:text-white"
            onClick={cancel}
          >
            <svg
              aria-hidden="true"
              class="w-5 h-5"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill-rule="evenodd"
                d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                clip-rule="evenodd"
              ></path>
            </svg>
            <span class="sr-only">Close modal</span>
          </button>
          <div class="py-6 px-6 lg:px-8">
            <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">
              Add Order
            </h3>
            <form class="space-y-6" action="#">
              <div>
                <label
                  for="email"
                  class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300"
                >
                  Price
                </label>
                <input
                  type="number"
                  min={0}
                  class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white"
                  placeholder="Name"
                  value={draftOrder().total_price}
                  onChange={(e) =>
                    setDraftOrder({
                      ...draftOrder(),
                      total_price: +(e.target as HTMLInputElement).value,
                    })
                  }
                />
                <label
                  for="suppliers"
                  class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-400"
                >
                  Select a supplier
                </label>
                {suppliers().length > 0 ? (
                  <select
                    id="suppliers"
                    class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white"
                    onChange={(e) => {
                      setDraftOrder({
                        ...draftOrder(),
                        supplier_id: (e.target as HTMLInputElement).value,
                      });
                    }}
                  >
                    <For each={suppliers()}>
                      {(supplier) => (
                        <option value={supplier.id}>{supplier.name}</option>
                      )}
                    </For>
                  </select>
                ) : (
                  <p>No suppliers found</p>
                )}
                <label
                  for="customers"
                  class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-400"
                >
                  Select a customer
                </label>
                {customers().length > 0 ? (
                  <select
                    id="customers"
                    class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white"
                    onChange={(e) => {
                      setDraftOrder({
                        ...draftOrder(),
                        customer_id: (e.target as HTMLInputElement).value,
                      });
                    }}
                  >
                    <For each={customers()}>
                      {(customer) => (
                        <option value={customer.id}>{customer.name}</option>
                      )}
                    </For>
                  </select>
                ) : (
                  <p>No customers found</p>
                )}
              </div>
              <button
                type="submit"
                onClick={(e) => {
                  e.preventDefault();
                  submit(draftOrder());
                }}
                class="w-full text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
              >
                Add
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};
export default AddOrderModal;
