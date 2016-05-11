DROP PROCEDURE IF EXISTS 	 `mopedmanage`.`getMopedBynameOrHphm1`  ;
DELIMITER $$
/*
select DISTINCT area_tb.area_name , moped_tb.moped_hphm ,type1_tb.dicword_wordname as typetype ,  color1_tb.dicword_wordname , 
	  owner_tb.owner_name , moped_tb.moped_id , tag_tb.tag_tagid ,moped_tb.moped_state , type2_tb.dicword_wordname as statestr
	  FROM owner_tb  JOIN moped_tb JOIN tag_tb   JOIN mopedowner_tb  
			ON moped_tb.moped_id = mopedowner_tb.moped_id AND  mopedowner_tb.owner_id = owner_tb.owner_id  
			JOIN mopedtag_tb ON mopedtag_tb.moped_id = moped_tb.moped_id AND mopedtag_tb.tag_id = tag_tb.tag_id  
			JOIN area_tb ON area_tb.area_id = moped_tb.area_id  
			JOIN  dicword_tb  AS type2_tb  ON  type2_tb.dicword_dictypeid = 8 AND tag_tb.tag_state = type2_tb.dicword_wordid  
			JOIN  dicword_tb  AS type1_tb  ON  type1_tb.dicword_dictypeid = 6 AND moped_tb.moped_type = type1_tb.dicword_wordid 
			JOIN   dicword_tb  AS color1_tb  ON   color1_tb.dicword_dictypeid = 7
			 AND moped_tb.moped_colorid = color1_tb.dicword_wordid  
			WHERE mopedtag_tb.mopedtag_state = 1 and moped_tb.moped_hphm like "%%%s%%" and owner_tb.owner_name like "%%%s%%" and (owner_tb.owner_state = 1) order by owner_tb.owner_id 
			
			*/
			
CREATE
    /*[DEFINER = { user | CURRENT_USER }]*/
    PROCEDURE `mopedmanage`.`getMopedBynameOrHphm1`(IN v_hphm VARCHAR(45) ,  IN v_ownername VARCHAR(45))
    /*LANGUAGE SQL
    | [NOT] DETERMINISTIC
    | { CONTAINS SQL | NO SQL | READS SQL DATA | MODIFIES SQL DATA }
    | SQL SECURITY { DEFINER | INVOKER }
    | COMMENT 'string'*/
    BEGIN
      
        DROP TABLE IF EXISTS tmp_out_tb;
        CREATE	TEMPORARY TABLE tmp_out_tb (   
			area_name  VARCHAR(45), 
			moped_hphm VARCHAR(45), 
			typetype VARCHAR(45),
			dicword_wordname VARCHAR(45),
			owner_name VARCHAR(45),
			moped_id  INTEGER ,
			tag_tagid VARCHAR(45),
			moped_state  INTEGER ,
			statestr VARCHAR(45)
			)  	;
		BEGIN 
			DECLARE s_Areaname VARCHAR(45);
			DECLARE s_Hphm VARCHAR(45);
			DECLARE s_Typetype VARCHAR(45);
			DECLARE s_Color VARCHAR(45);
			DECLARE s_Ownername VARCHAR(45);
			DECLARE i_Moped_id INT(10) ;
			DECLARE i_Moped_state INT(10);
			DECLARE s_Tag_tagid VARCHAR(45);
			DECLARE s_tag_state VARCHAR(45);
			
			DECLARE i_moped_type INT(10);
			DECLARE i_areaid INT(10);
			DECLARE i_colorid INT(10);
			DECLARE i_ownerid INT(10);
			DECLARE i_Tag_tagid INT(10);
			DECLARE i_tag_state INT(10);
		   
			DECLARE done INT DEFAULT 0;
			
		   
		   
			DECLARE p_moped_tb CURSOR FOR SELECT DISTINCT area_tb.area_name , moped_tb.moped_hphm ,moped_tb.moped_type ,  moped_tb.moped_colorid ,  owner_tb.owner_name , moped_tb.moped_id , tag_tb.tag_tagid ,moped_tb.moped_state , tag_tb.tag_state
	        FROM owner_tb  JOIN moped_tb JOIN tag_tb   JOIN mopedowner_tb  
			ON moped_tb.moped_id = mopedowner_tb.moped_id AND  mopedowner_tb.owner_id = owner_tb.owner_id
			JOIN mopedtag_tb ON mopedtag_tb.moped_id = moped_tb.moped_id AND mopedtag_tb.tag_id = tag_tb.tag_id  
			JOIN area_tb ON area_tb.area_id = moped_tb.area_id
			WHERE mopedtag_tb.mopedtag_state = 1 AND moped_tb.moped_hphm LIKE  CONCAT('%',v_hphm,'%') AND owner_tb.owner_name LIKE CONCAT('%',v_ownername,'%') AND (owner_tb.owner_state = 1) ORDER BY owner_tb.owner_id ;
  
			
			DECLARE CONTINUE HANDLER FOR NOT FOUND SET done=1;
			
			OPEN p_moped_tb  ;
			cursor_loop_moped : LOOP	
						
				FETCH p_moped_tb INTO s_Areaname ,   s_Hphm , i_moped_type , i_colorid , s_Ownername , i_Moped_id , s_Tag_tagid , i_Moped_state , i_tag_state ;
				BEGIN 
					SELECT dicword_wordname INTO s_Color FROM dicword_tb WHERE  dicword_dictypeid = 7 AND dicword_wordid =  i_colorid  AND dicword_state = 1;
					SELECT dicword_wordname INTO s_Typetype FROM dicword_tb WHERE  dicword_dictypeid = 6 AND dicword_wordid =  i_moped_type AND dicword_state = 1;
					SELECT dicword_wordname INTO s_tag_state FROM dicword_tb WHERE  dicword_dictypeid = 8 AND dicword_wordid =  i_tag_state AND dicword_state = 1;
					INSERT INTO tmp_out_tb VALUES(s_Areaname ,s_Hphm ,s_Typetype , s_Color ,s_Ownername ,i_moped_id ,s_Tag_tagid ,i_Moped_state , s_tag_state ) ;
										
					
				END;
                IF done=1 THEN
					
					LEAVE cursor_loop_moped;
				END IF;				
			END LOOP cursor_loop_moped ;	
			CLOSE p_moped_tb ;
		END;
	
	SELECT DISTINCT * FROM tmp_out_tb ;
	 
    END$$

DELIMITER ;