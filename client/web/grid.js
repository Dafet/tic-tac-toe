class Grid {
    _elems = ["-", "-", "-", "-", "-", "-", "-", "-", "-"];

    get elems() {
        return this._elems;
    }

    set elems(value) {
        this._elems = value;
    }

    fillCell(id, mark) {
        if (!this.validCellPlacement(id)) {
            return false
        }

        this.elems[id] = mark;

        return true
    }

    validCellPlacement(id) {
        if (id < 0 || id > 8) {
            console.log("[error] cell ID must be in range 0 - 8")
            return false
        }

        if (this.cellOccupied(id)) {
            console.log("[error] cell %s is occupied", id)
            return false
        }

        return true
    }

    cellOccupied(id) {
        if (this.elems[id] != "-") {
            return true
        }

        return false
    }

    draw() {
        for (let i = 0; i < 9; i++) {
            const element = this.elems[i];
            const cellElem = document.getElementById("cell-" + i);
            cellElem.textContent = element;
        }
    }

    reset() {
        this.elems = ["-", "-", "-", "-", "-", "-", "-", "-", "-"];
    }
}