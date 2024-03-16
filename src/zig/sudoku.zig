const std = @import("std");

const GRID_SIZE: u64 = 9;

// TODO(ElodinLaarz): allow variable sized sudoku grids.
const sudokuGrid = struct {
    grid_size: u64,
    cells: u8[GRID_SIZE][GRID_SIZE],
    // These feel inefficient...
    valid_values_by_cell_index: std.AutoHashmap(u8, [GRID_SIZE]bool),
    valid_cell_index_by_value: std.AutoHashmap(u8, [81]bool),

    pub fn calculate_valid_cells(self: sudokuGrid) void {
        for (0..GRID_SIZE) |i| {
            for (0..GRID_SIZE) |j| {
                if (self.cells[i][j] == 0) {
                    for (0..GRID_SIZE) |k| {
                        self.valid_values_by_cell_index[i][j][k] = true;
                    }
                } else {
                    self.valid_values_by_cell_index[i][j][self.cells[i][j]] = false;
                }
            }
        }
    }
};

fn createInit(alloc: std.mem.Allocator, comptime T: type, props: anytype) !*T {
    const new = try alloc.create(T);
    new.* = props;
    return new;
}

fn printSudoku(grid: *sudokuGrid) void {
    for (grid.cells) |row| {
        for (row) |cell| {
            if (cell == 0) {
                var valid_cells = "";
                for (0..grid.grid_size) |i| {
                    if (grid.valid_values_by_cell_index[i][cell]) {
                        valid_cells += std.fmt.sprintf("%d.", i);
                    }
                }
            }
            std.debug.print_u8(cell);
            std.debug.print_str(". ");
        }
        std.debug.print_str("\n");
    }
}

// fn hasNeighbor(grid: *sudokuGrid, row: u8, col: u8, num: u8) bool {
//     var i = 0;
//     while (i < 9) : (i += 1) {
//         if (grid.cells[row][i] == num) {
//             return true;
//         }
//         if (grid.cells[i][col] == num) {
//             return true;
//         }
//     }
//     var startRow = row - row % 3;
//     var startCol = col - col % 3;
//     var r = startRow;
//     while (r < startRow + 3) : (r += 1) {
//         var c = startCol;
//         while (c < startCol + 3) : (c += 1) {
//             if (grid.cells[r][c] == num) {
//                 return true;
//             }
//         }
//     }
//     return false;
// }

// fn isNonZeroSquare(grid: *sudokuGrid) bool {
//     if (grid.cells.len() == 0) {
//         return false;
//     }
//     var row_len = grid.cells[0].len();
//     const row_and_col_len: bool = row_len == grid.cells.len();

//     var equal_rows: bool = true;
//     for (grid.cells) |row| {
//         if (row.len() != row_len) {
//             equal_rows = false;
//             break;
//         }
//     }
//     return row_and_col_len and equal_rows;
// }

fn randomAdd(grid: *sudokuGrid) [][2]u8 {
    var num_to_add: u8 = 1;
    const side_length: u64 = grid.cells.len();
    _ = side_length;
    const grid_size: u64 = grid.grid_size;

    while (num_to_add <= grid_size) : (num_to_add += 1) {
        const row = std.rand.intn(grid_size);
        const col = std.rand.intn(grid_size);
        if (grid.cells[row][col] == 0) {
            grid.cells[row][col] = num_to_add;
        } else {
            continue;
        }
    }
}

fn removeAtIndices(grid: *sudokuGrid, indices: [][2]u8) void {
    for (indices) |index| {
        grid.cells[index[0]][index[1]] = 0;
    }
    grid.calculate_valid_cells();
}

// TODO(ElodinLaarz): allow variable sized sudoku grids.
fn createSudoku() sudokuGrid {
    // var side_length: u64 = box_size * box_size;
    // var grid_size: u64 = side_length * side_length;
    var grid: sudokuGrid = try createInit(sudokuGrid, .{
        .grid_size = GRID_SIZE,
        .cells = [GRID_SIZE]u8{},
        .valid_values_by_cell_index = std.AutoHashmap(u8, [GRID_SIZE]bool),
        .valid_cell_index_by_value = std.AutoHashmap(u8, [81]bool),
    });

    const blank_cells: u8 = 81;
    var consecutive_failures: u8 = 0;
    const max_failures: u8 = 2;
    const added_indices: [][9]u8 = [][9]u8{};

    while (blank_cells > 0) {
        var new_indices: [][2]u8 = randomAdd(&grid);
        if (new_indices.len() == 0) {
            consecutive_failures += 1;
            if (consecutive_failures > max_failures) {
                break;
            }
            if (blank_cells < grid.grid_size) {
                removeAtIndices(&grid, added_indices);
            }
        } else {}
    }
    return grid;
}

pub fn main() void {
    var grid: sudokuGrid = createSudoku();
    printSudoku(&grid);
}
