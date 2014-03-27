<?php
class User_db extends CI_Model {

	private $table	=	'category';

	/**
	 * inserts a new entry in category table
	 *
	 * @author Anil
	 * @param string $cat_name
	 * @param integer $parent_id
	 * @return integer last inserted id
	 */
	function insertNewCategory($cat_name,$parent_id = 0){

		$check_duplicate = $this->getCatId($cat_name, $parent_id);

		if($check_duplicate > 0){
			log_message('error','Its a duplicate entry category with name:'.$cat_name.' under parent : '.$parent_id);
			return $check_duplicate;
		}

		$this->db->insert($this->table, array(
			'cat_name' => $cat_name
			,'parent_id' => $parent_id
		));

		if ($this->db->_error_message()){
			$err_msg = $this->db->_error_message();
			$err_num = $this->db->_error_number();
			log_message('error',$err_num.' : '.$err_msg.'. It happened during inserting new category with name:'.$cat_name.' under parent : '.$parent_id);
			return 0;
		}else{
			$new_cat_id = $this->db->insert_id();
			return $new_cat_id;
		}
	}

	function getCatId($cat_name,$parent_id = 0){
		$result = $this->db->get_where($this->table,array('cat_name' => $cat_name, 'parent_id' => $parent_id));

		if($result->num_rows() > 0){
			//ideally we should get only one row
			$row	= $result->row();
			return $row->cat_id;
		}else
			return 0;
	}

	function getParentsUptoCertainLevel($cid, $level = 2){
		if(!$cid)
			return array();

		$list = array();

		for( ;$level > 0 ; $level--){
			$result = $this->db->select('cat_id,cat_name,parent_id')
			->from($this->table)
			->where(array('cat_id'=>$cid,'deleted !=' => 1))
			->get();

			$row = $result->row();

			$list[] = array(
				'cat_id' => $row->cat_id,
				'cat_name' => $row->cat_name,
				'cat_par' => $row->parent_id,
			);

			$cid = $row->parent_id;
			if($cid == 0)
				break;
		}

		return array_reverse($list);
	}
}